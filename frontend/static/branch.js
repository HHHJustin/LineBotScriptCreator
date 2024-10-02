import { getCurrentIDAndTypeFromURL } from './utils.js';

export async function addBranch() {
    const { currentID } = getCurrentIDAndTypeFromURL(); 
    const currentIDInt = parseInt(currentID, 10);
    const path = "/nodes/create/branch";  
    const nodeType = document.getElementById("TypeSelect").value;
    if (!nodeType) {
        alert("請選擇一個分支類型");
        return;
    }
    const data = {
        CurrentNodeID: currentIDInt,
        NewNodeType: nodeType
    };
    
    try {
        const response = await fetch(path, {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify(data),
        });

        if (!response.ok) {
            throw new Error('Network response was not ok');
        }

        const result = await response.json();
        console.log('Branch created successfully:', result);
        window.location.reload();
    } catch (error) {
        console.error('Error creating branch:', error);
        alert('Error creating branch: ' + error.message);  
    }
}
