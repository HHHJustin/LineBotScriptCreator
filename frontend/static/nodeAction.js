import { getCurrentIDAndTypeFromURL } from './utils.js';
let myDiagramInstance = null;

export function initialize(diagram) {
    myDiagramInstance = diagram;
}

export function showContextMenu(event, nodeData) {
    let menu = document.getElementById('contextMenu');
    if (!menu) {
        menu = document.createElement('div');
        menu.id = 'contextMenu';
        menu.style.position = 'absolute';
        menu.style.backgroundColor = '#fff';
        menu.style.border = '1px solid #ccc';
        menu.style.padding = '10px';
        menu.style.boxShadow = '0 2px 10px rgba(0,0,0,0.2)';
        menu.style.zIndex = 1000;
        menu.style.display = 'none';
        document.body.appendChild(menu);
    }

    menu.innerHTML = '';
    const addButton = document.createElement('button');
    addButton.textContent = 'Add';
    addButton.onclick = async () => {
        await createNodeForParent(nodeData);
        hideMenu();
    };
    menu.appendChild(addButton);

    const removeButton = document.createElement('button');
    removeButton.textContent = 'Remove';
    removeButton.onclick = () => {
        removeNode(nodeData);
        hideMenu();
    };
    menu.appendChild(removeButton);
    const menuX = event.clientX;
    const menuY = event.clientY;

    menu.style.left = `${menuX}px`;
    menu.style.top = `${menuY}px`;
    menu.style.display = 'block';
    document.addEventListener('click', hideMenuOutside);
}

function hideMenu() {
    const menu = document.getElementById('contextMenu');
    if (menu) {
        menu.style.display = 'none';
    }
    document.removeEventListener('click', hideMenuOutside);
}

function hideMenuOutside(event) {
    const menu = document.getElementById('contextMenu');
    if (menu && !menu.contains(event.target)) {
        menu.style.display = 'none';
        document.removeEventListener('click', hideMenuOutside);
    }
}

async function createNodeForParent(parentNodeData) {
    const typeId = prompt("請輸入新節點的類型ID"); 
    await createNode(typeId, parentNodeData);
}

function removeNode(nodeData) {
    myDiagramInstance.model.remove(nodeData);
}

export function goToNode(nodeID) {
    const url = `/nodes/type?currentNodeID=${nodeID}`;
    fetch(url)
        .then(response => response.json())
        .then(data => {
            const nodeType = data.nodeType;

            if (nodeType === "Message") {
                window.location.href = `/nodes/get/${nodeID}/Message`;
            } else if (nodeType === "QuickReply") {
                window.location.href = `/nodes/get/${nodeID}/QuickReply`;
            }  else if (nodeType === "KeywordDecision") {
                window.location.href = `/nodes/get/${nodeID}/KeywordDecision`;
            }  else if (nodeType === "TagDecision") {
                window.location.href = `/nodes/get/${nodeID}/TagDecision`;
            }  else if (nodeType === "TagOperation") {
                window.location.href = `/nodes/get/${nodeID}/TagOperation`;
            }  else if (nodeType === "Random") {
                window.location.href = `/nodes/get/${nodeID}/Random`;
            }  else if (nodeType === "FirstStep") {
                window.location.href = `/firststep/read`;
            }  else {
                console.error("Unsupported node type:", nodeType);
            }
        })
        .catch(error => {
            console.error('Error fetching node type:', error);
        });
}

export function makeTitleEditable(td, nodeID) {
    const originalTitle = td.innerText;
    const input = document.createElement('input');
    input.type = 'text';
    input.value = originalTitle;
    input.onkeydown = async function(event) {
        if (event.key === 'Enter') {
            const newTitle = input.value;
            try {
                const response = await fetch('/nodes/title', {
                    method: 'POST',
                    headers: {
                        'Content-Type': 'application/json',
                    },
                    body: JSON.stringify({
                        currentNodeID: nodeID,
                        newTitle: newTitle,
                    }),
                });

                if (response.ok) {
                    const jsonResponse = await response.json();
                    console.log('Node title updated:', jsonResponse);
                    td.innerText = newTitle; 
                } else {
                    console.error('Error updating node title');
                }
            } catch (error) {
                console.error('Error:', error);
            }
        }
    };

    td.innerHTML = '';
    td.appendChild(input);
    input.focus();
}


let draggedRow = null;

export function dragStart(event) {
    draggedRow = event.target.closest('tr'); 
    event.dataTransfer.effectAllowed = "move";
}

export function allowDrop(event) {
    event.preventDefault();  
}

export function drop(event) {
    event.preventDefault();
    const { currentID } = getCurrentIDAndTypeFromURL(); 
    const currentIDInt = parseInt(currentID, 10);
    const dropTargetRow = event.target.closest('tr'); 
    if (dropTargetRow && draggedRow && draggedRow !== dropTargetRow) {
        const tbody = document.getElementById("messageTableBody");
        const rows = Array.from(tbody.querySelectorAll("tr[draggable='true']"));
        const draggedMessageIndex = rows.indexOf(draggedRow); 
        const newIndex = rows.indexOf(dropTargetRow);

        var data = {
            currentNodeID: currentIDInt,
            draggedMessageIndex: draggedMessageIndex, 
            newIndex: newIndex, 
        };

        console.log(data);
        fetch('/messages/updateorder', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify(data),
        })
        .then(response => {
            if (!response.ok) {
                throw new Error('Failed to update message order');
            }
            return response.json();
        })
        .then(data => {
            console.log('Order updated successfully:', data);
            window.location.reload();
        })
        .catch(error => {
            console.error('Error updating message order:', error);
        });
    }
}