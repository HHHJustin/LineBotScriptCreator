import { getCurrentIDAndTypeFromURL } from './utils.js';

export async function addQuickReply() {
    const { currentID } = getCurrentIDAndTypeFromURL(); 
    const currentIDInt = parseInt(currentID, 10);
    const buttonName = document.getElementById("newButtonName").value;
    const reply = document.getElementById("newReply").value;
    console.log(currentIDInt)
    if (!buttonName || !reply) {
        alert("Please fill in both Button Name and Reply fields.");
        return;
    }

    const data = {
        currentNodeID: currentIDInt,
        buttonName: buttonName,
        reply: reply
    };
    
    try {
        const response = await fetch('/quickreplies/create', {
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
        console.log('Quick reply added successfully:', result);
        window.location.reload();
    } catch (error) {
        console.error('Error adding quick reply:', error);
        alert('Error adding quick reply: ' + error.message);
    }
}

export async function deleteQuickReply(event, quickReplyID) {
    event.preventDefault();  
    const { currentID } = getCurrentIDAndTypeFromURL(); 
    if (!currentID) {
        console.error("Node ID is not found in the URL.");
        return;
    }
    const currentIDInt = parseInt(currentID, 10);
    if (isNaN(currentIDInt)) {
        console.error("Node ID is not a valid number.");
        return;
    }
    const data = {
        QuickReplyID: quickReplyID,
        currentNodeID: currentIDInt
    };
    fetch('/quickreplies/delete', {
        method: 'POST',  
        headers: {
            'Content-Type': 'application/json',  
        },
        body: JSON.stringify(data)  
    })
    .then(response => {
        if (!response.ok) {
            throw new Error("Network response was not ok");
        }
        return response.json();  
    })
    .then(data => {
        console.log('Success:', data);
        window.location.reload();  
    })
    .catch(error => {
        console.error('Error:', error);
    });
}

export function makeQuickReplyEditable(td, quickReplyID, field) {
    const originalContent = td.innerText;
    const input = document.createElement('input');
    input.type = 'text';
    input.value = originalContent;
    input.onkeydown = async function(event) {
        if (event.key === 'Enter') {
            const updatedContent = input.value;

            const data = {
                quickReplyID: quickReplyID,
                field: field,
                value: updatedContent
            };

            try {
                const response = await fetch('/quickreplies/update', {
                    method: 'POST',
                    headers: {
                        'Content-Type': 'application/json',
                    },
                    body: JSON.stringify(data),
                });

                if (response.ok) {
                    const jsonResponse = await response.json();
                    console.log('Quick Reply updated:', jsonResponse);
                    td.innerText = updatedContent; 
                } else {
                    console.error('Error updating Quick Reply');
                    td.innerText = originalContent; 
                }
            } catch (error) {
                console.error('Error:', error);
                td.innerText = originalContent; 
            }
        }
    };

    td.innerHTML = '';
    td.appendChild(input);
    input.focus();
}