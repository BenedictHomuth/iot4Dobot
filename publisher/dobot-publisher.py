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

# Retrieve current robot position
roboPos = dType.GetPose(api)

# Generate json ready data
posData = buildJsonRequest(roboPos)
print("Generated Positional Data:")
print(posData)

# Send data to backend -> which then sends it to NATS
response = requests.post('https://159.69.125.114:8080/event', json = posData, verify = False)

print("Web Request send: Status " + str(response.status_code) + response.text)