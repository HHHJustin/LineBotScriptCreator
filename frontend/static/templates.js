var $ = go.GraphObject.make; 
var myDiagram; 

export function setupNodeTemplate(diagram) {
    myDiagram = diagram;
    myDiagram.nodeTemplate =
        $(go.Node, "Auto",
            {
                movable: true, 
                click: null,
                contextMenu: createContextMenu()
            },
            new go.Binding("location", "loc", go.Point.parse).makeTwoWay(go.Point.stringify),
            $(go.Shape, "RoundedRectangle",
                { strokeWidth: 0 },
                new go.Binding("fill", "color")
            ),
            $(go.TextBlock,
                { margin: 8, font: "bold 14px sans-serif" },
                new go.Binding("text", "text")
            )
        );
    myDiagram.linkTemplate =
    $(go.Link,
        $(go.Shape), 
        $(go.Shape, { toArrow: "Standard" }),  
        {  
            contextMenu:
                $(go.Adornment, "Vertical",
                    $("ContextMenuButton",
                        $(go.TextBlock, "Delete Link", { margin: 10, font: "bold 14px sans-serif" }),
                        {
                            click: function(e, obj) {
                                deleteLinkHandler(e, obj);
                            }
                        }
                    )
                )
        }
    );
}

function createContextMenu() {
    console.log('Creating context menu');
    return $(go.Adornment, "Spot",
        $(go.Panel, "Vertical",
            {
                alignment: go.Spot.Center,
                defaultAlignment: go.Spot.Left,
                margin: 5
            },
            createAddNextNodeMenuItem(),     
            createAddPreviousNodeMenuItem(), 
            createDeleteNodeMenuItem(),
            createAddLinkFromMenuItem(),
            createAddLinkToMenuItem(),
            createAddBranchMenuItem(),
            createEditMenuItem()
        )
    );
}

/* Create Node */
function createAddNextNodeMenuItem() {
    return $("ContextMenuButton",
        $(go.Panel, "Horizontal",
            $(go.TextBlock, "Add Next Node",
                {
                    font: "bold 16px sans-serif",
                    margin: new go.Margin(10, 25, 10, 25),
                }
            ),
        ),
        {
            width: 200,
            height: 50,
            click: function(e, obj) { addNodeHandler(e, obj, "next"); }
        }
    );
}

function createAddPreviousNodeMenuItem() {
    return $("ContextMenuButton",
        $(go.Panel, "Horizontal",
            $(go.TextBlock, "Add Previous Node",
                {
                    font: "bold 16px sans-serif",
                    margin: new go.Margin(10, 25, 10, 25),
                    alignment: go.Spot.Left
                }
            ),
        ),
        {
            width: 200,
            height: 50,
            click: function(e, obj) { addNodeHandler(e, obj, "previous") }
        }
    );
}

function addNodeHandler(e, obj, type) {
    console.log("Clicked: Add " + type + " Node");
    var node = obj.part.adornedPart;  
    var currentNodeID = node.data.key; 
    var path = '/nodes/create/' + type;
    showCenterDialog(path, currentNodeID);
}

/* Delete Node */
function createDeleteNodeMenuItem() {
    return $("ContextMenuButton",
        $(go.Panel, "Horizontal",
            $(go.TextBlock, "Delete Node",
                {
                    font: "bold 16px sans-serif",
                    margin: new go.Margin(10, 25, 10, 25),
                    alignment: go.Spot.Left
                }
            ),
        ),
        {
            width: 200,
            height: 50,
            click: function(e, obj) { deleteNodeHandler(e, obj); }
        }
    );
}

function deleteNodeHandler(e, obj) {
    var node = obj.part.adornedPart;
    var currentNodeID = node.data.key;
    var data = {
        currentNodeID: currentNodeID
    };
    var path = '/nodes/delete';
    sendDataToServer(data, path);
}

/* Add Link */
function createAddLinkFromMenuItem() {
    return $("ContextMenuButton",
        $(go.Panel, "Horizontal",
            $(go.TextBlock, "Add Link (From)",
                {
                    font: "bold 16px sans-serif",
                    margin: new go.Margin(10, 25, 10, 25),
                    alignment: go.Spot.Left
                }
            ),
        ),
        {
            width: 200,
            height: 50,
            click: function (e, obj) {
                var node = obj.part.adornedPart;  
                myDiagram.startTransaction("setLinkFrom");
                myDiagram.model.setDataProperty(node.data, "isLinkFromSelected", true); 
                myDiagram.commitTransaction("setLinkFrom");
                console.log("Set Link From Node:", node.data.key);
            }
        }
    );
}

function createAddLinkToMenuItem() {
    return $("ContextMenuButton",
        $(go.Panel, "Horizontal",
            $(go.TextBlock, "Add Link (To)",
                {
                    font: "bold 16px sans-serif",
                    margin: new go.Margin(10, 25, 10, 25),
                    alignment: go.Spot.Left
                }
            )
        ),
        {
            width: 200,
            height: 50,
            click: function (e, obj) {
                var toNode = obj.part.adornedPart;  
                var fromNode = findFromNode();  
                if (!fromNode) {
                    alert("Please select 'Add Link (From)' first.");
                    return;
                }
                addLinkHandler(fromNode, toNode);  
            }
        }
    );
}

function addLinkHandler(fromNode, toNode) {
    var data = {
        fromNodeID: fromNode.data.key, 
        toNodeID: toNode.data.key       
    };
    var path = '/links/create';
    sendDataToServer(data, path);  
    resetLinkFromStatus();
}

function findFromNode() {
    var nodes = myDiagram.model.nodeDataArray;
    for (var i = 0; i < nodes.length; i++) {
        if (nodes[i].isLinkFromSelected) {
            console.log("Found Link From Node:", nodes[i]);  
            return myDiagram.findNodeForData(nodes[i]);
        }
    }
    console.log("No Link From Node found.");
    return null;  
}

function resetLinkFromStatus() {
    myDiagram.startTransaction("resetLinkFromStatus");
    myDiagram.model.nodeDataArray.forEach(function (nodeData) {
        myDiagram.model.setDataProperty(nodeData, "isLinkFromSelected", false);
    });
    myDiagram.commitTransaction("resetLinkFromStatus");
}

/* Delete Link */
function deleteLinkHandler(e, obj){
    var link = obj.part; 
    console.log(link)
    if (link !== null) {
        var fromNodeKey = link.data.from;
        var toNodeKey = link.data.to;
        var fromNodeData = myDiagram.model.findNodeDataForKey(fromNodeKey);
        var toNodeData = myDiagram.model.findNodeDataForKey(toNodeKey);
        var data = {
            fromNodeID: fromNodeData.key,
            toNodeID: toNodeData.key
        };
        var path = '/links/delete';
        sendDataToServer(data, path);
    }
}

/* Add Branch */
function createAddBranchMenuItem() {
    return $("ContextMenuButton",
        $(go.Panel, "Horizontal",
            $(go.TextBlock, "Add Branch",
                {
                    font: "bold 16px sans-serif",
                    margin: new go.Margin(10, 25, 10, 25),
                    alignment: go.Spot.Left
                }
            ),
        ),
        {
            width: 200,
            height: 50,
            click: function(e, obj) { addbranchHandler(e, obj); }
        }
    );
}

function addbranchHandler(e, obj) {
    var node = obj.part.adornedPart;
    var currentNodeID = node.data.key;
    var path = `/nodes/type?currentNodeID=${currentNodeID}`;
    fetch(path, {
        method: 'GET',
        headers: {
            'Content-Type': 'application/json'
        }
    })
    .then(response => response.json())
    .then(data => {
        console.log('Success:', data);
        const nodeType = data.nodeType;
        const nodeID = data.nodeID;
        if (nodeType === "Message" || nodeType === "QuickReply" || nodeType === "TagOperation" || nodeType === "FirstStep") {
            alert(`Unsupported operation: Node type ${nodeType} is not allowed.`);
        } else if (nodeType === "KeywordDecision") {
            var path = '/nodes/create/branch';
            showCenterDialog(path, currentNodeID);
        } else if (nodeType === "TagDecision") {
            window.location.href = `/nodes/get/${nodeID}/TagDecision`;
        } else if (nodeType === "Random") {
            window.location.href = `/nodes/get/${nodeID}/Random`;
        } else {
            console.log("Unsupported node type:", nodeType);
        }
    })
    .catch((error) => {
        console.error('Error:', error);
    });
}

/* Edit */
function createEditMenuItem() {
    return $("ContextMenuButton",
        $(go.Panel, "Horizontal",
            $(go.TextBlock, "Edit",
                {
                    font: "bold 16px sans-serif",
                    margin: new go.Margin(10, 25, 10, 25),
                    alignment: go.Spot.Left
                }
            ),
        ),
        {
            width: 200,
            height: 50,
            click: function(e, obj) { editHandler(e, obj); }
        }
    );
}

function editHandler(e, obj) {
    var node = obj.part.adornedPart;
    var currentNodeID = node.data.key;
    var path = `/nodes/type?currentNodeID=${currentNodeID}`;
    fetch(path, {
        method: 'GET',
        headers: {
            'Content-Type': 'application/json'
        }
    })
    .then(response => response.json())
    .then(data => {
        console.log('Success:', data);
        const nodeType = data.nodeType;
        const nodeID = data.nodeID;
        if (nodeType === "Message") {
            window.location.href = `/nodes/get/${nodeID}/Message`;
        } else if (nodeType === "QuickReply") {
            window.location.href = `/nodes/get/${nodeID}/QuickReply`;
        } else if (nodeType === "KeywordDecision") {
            window.location.href = `/nodes/get/${nodeID}/KeywordDecision`;
        }  else if (nodeType === "TagDecision") {
            window.location.href = `/nodes/get/${nodeID}/TagDecision`;
        }  else if (nodeType === "TagOperation") {
            window.location.href = `/nodes/get/${nodeID}/TagOperation`;
        }  else if (nodeType === "Random") {
            window.location.href = `/nodes/get/${nodeID}/Random`;
        }  else if (nodeType === "FirstStep") {
            window.location.href = `/firststep/read`;
        }else {
            console.log("Unsupported node type:", nodeType);
        }
    })
    .catch((error) => {
        console.error('Error:', error);
    });
}

function sendDataToServer(data, path) {
    fetch(path, {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json'
        },
        body: JSON.stringify(data)
    })
    .then(response => {
        if (!response.ok) {
            throw new Error("Network response was not ok");
        }
        return response.text();  
    })
    .then(text => {
        if (text) {
            try {
                const jsonData = JSON.parse(text);  
                console.log('Success:', jsonData);
                window.location.reload();
            } catch (error) {
                console.error('Error parsing JSON:', error);
            }
        } else {
            console.log('No JSON response returned from server.');
        }
    })
    .catch((error) => {
        console.error('Error:', error);
    });
}

function showCenterDialog(path, currentNodeID) {
    var dialog = document.createElement("div");
    dialog.style.position = "fixed";
    dialog.style.left = "50%";
    dialog.style.top = "50%";
    dialog.style.transform = "translate(-50%, -50%)";
    dialog.style.zIndex = "1000"; 
    dialog.style.width = "300px";  
    dialog.style.padding = "10px";  
    dialog.style.backgroundColor = "transparent"; 
    
    var options = ["Message", "QuickReply", "KeywordDecision", "TagDecision", "TagOperation", "Random"];
    options.forEach(function(option) {
        var button = document.createElement("button");
        button.innerHTML = option;
        button.style.font = "bold 16px sans-serif";
        button.style.display = "block";
        button.style.width = "100%";
        button.style.margin = "10px 0";
        button.style.width = "200px";              
        button.style.height = "50px";                
        button.style.margin = "10px 25px";          
        button.onclick = function() {
            var data = {
                currentNodeID: currentNodeID, 
                newNodeType: option
            };
            sendDataToServer(data, path);
            document.body.removeChild(dialog); 
        };
        dialog.appendChild(button);
    });
    var cancelButton = document.createElement("button");
    cancelButton.innerHTML = "Cancel";
    cancelButton.style.font = "bold 16px sans-serif";
    cancelButton.style.display = "block";
    cancelButton.style.width = "200px";  
    cancelButton.style.height = "50px";  
    cancelButton.style.margin = "10px 25px"; 
    cancelButton.style.backgroundColor = "red";  
    cancelButton.style.color = "black";  
    cancelButton.onclick = function() {
        document.body.removeChild(dialog);
    };
    dialog.appendChild(cancelButton); 
    document.body.appendChild(dialog);
}