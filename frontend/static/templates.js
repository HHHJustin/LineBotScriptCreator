var $ = go.GraphObject.make; 
var myDiagram; 

export function setupNodeTemplate(diagram) {
    myDiagram = diagram;
    myDiagram.nodeTemplate =
        $(go.Node, "Auto",
            {
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
}

function createContextMenu() {
    return $(go.Adornment, "Spot",
        $(go.Panel, "Vertical",
            {
                alignment: go.Spot.Center,
                defaultAlignment: go.Spot.Left,
                margin: 5
            },
            createAddNextMenuItem(),
            createAddPreviousMenuItem(),
            createAddBranchMenuItem(),
            createRemoveMenuItem(),
            createAddLinkMenuItem(),
            createEditMenuItem()
        )
    );
}

function createAddNextMenuItem() {
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

function createAddPreviousMenuItem() {
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

function createRemoveMenuItem() {
    return $("ContextMenuButton",
        $(go.Panel, "Horizontal",
            $(go.TextBlock, "Remove Node",
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
            click: function(e, obj) { removeNodeHandler(e, obj); }
        }
    );
}

function createAddLinkMenuItem() {
    return $("ContextMenuButton",
        $(go.Panel, "Horizontal",
            $(go.TextBlock, "Remove Node",
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
            click: function(e, obj) { addLinkHandler(e, obj); }
        }
    );
}

function createEditMenuItem() {
    return $("ContextMenuButton",
        $(go.Panel, "Horizontal",
            $(go.TextBlock, "Remove Node",
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

function addNodeHandler(e, obj, type) {
    console.log("Clicked: Add " + type + " Node");
    showCenterDialog(type, obj);
}

function showCenterDialog(type, obj) {
    var dialog = document.createElement("div");
    dialog.style.position = "fixed";
    dialog.style.left = "50%";
    dialog.style.top = "50%";
    dialog.style.transform = "translate(-50%, -50%)";
    dialog.style.zIndex = "1000"; 
    dialog.style.width = "300px";  
    dialog.style.padding = "10px";  
    dialog.style.backgroundColor = "transparent"; 
    var node = obj.part.adornedPart;  
    var currentNodeID = node.data.key; 
    var options = ["Message", "QuickReply", "Keyword Decision", "Tag Decision", "Add Tag", "Cancel Tag", "Random"];
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
            var path = '/nodes/create/' + type;
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


function addbranchHandler(e, obj) {

}

function removeNodeHandler(e, obj) {

}

function addLinkHandler(e, obj) {

}

function editHandler(e, obj) {

}

function sendDataToServer(data, path) {
    fetch(path, {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json'
        },
        body: JSON.stringify(data)  
    })
    .then(response => response.json())  
    .then(data => {
        console.log('Success:', data);
        window.location.reload();
    })
    .catch((error) => {
        console.error('Error:', error);
    });
}
