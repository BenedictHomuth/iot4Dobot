import threading
import time
import requests

def buildJsonRequest(poseData):
	# roboEvent is the blueprint for the json data
	roboEvent = {
			 "x": 0, # cartesian x
			 "y": 0, # cartesian y
			 "z": 0, # cartesian z
			 "r": 0, # cartesian r -> end effector rotation
			 "jointAngles": [ 0, # joint rotation j1 -> rearArm
							 0, # joint rotation j2 -> foreArm
							 0, # joint rotation j3 -> z-axis position
							 0 # joint rotation j4 -> end effector
						 ]
			}
	
	# Assign cartesian values
	roboEvent["x"] = poseData[0]
	roboEvent["y"] = poseData[1]
	roboEvent["z"] = poseData[2]
	roboEvent["r"] = poseData[3]
	
	# Assign joint rotations
	roboEvent["jointAngles"][0] = poseData[4]
	roboEvent["jointAngles"][1] = poseData[5]
	roboEvent["jointAngles"][2] = poseData[6]
	roboEvent["jointAngles"][3] = poseData[7]
	return roboEvent

def getPose():
	while True:
		roboPos = dType.GetPose(api)
		posData = buildJsonRequest(roboPos)
		
		# Send data to backend -> which then sends it to NATS
		response = requests.post('https://159.69.125.114:8080/event', json = posData, verify = False)
		#print("Web Request send: Status " + str(response.status_code) + response.text)
		
		# get the current robot arm position every 100ms
		time.sleep(0.1)
		
def primaryProgram():
	while True:
		dType.SetPTPCmdSync(api, 0, 200, 0, 100, 0, 1)
		dType.SetPTPCmdSync(api, 0, 200, 100, 100, 0, 1)
		dType.SetPTPCmdSync(api, 0, 400, 0, 230, 0, 1)
		

# Set Arm Parameters
dType.SetPTPCommonParams(api, 30, 30, 0)
dType.SetArmOrientation(api, 1, 1)

# create a thread for the getPose function
pose_thread = threading.Thread(target=getPose)

# This kills the thread, as soon as the primary programms finishes / cancelles
pose_thread.daemon = True

# Start the thread
pose_thread.start()

# Run the primary program in the main thread
primaryProgram()