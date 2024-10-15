export async function submitFirstStep() {
    var firstStepType = document.getElementById("firstStepTypeSelect").value;

    if (!firstStepType) {
        alert("請填寫完整的資料");
        return;
    }

    var data = {
        firstStepType: firstStepType
    };

    try {
        const response = await fetch('/nodes/create/firststep', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify(data)
        });

        if (!response.ok) {
            throw new Error('Network response was not ok');
        }

        const result = await response.json();
        console.log('Success:', result);
        window.location.reload();

    } catch (error) {
        console.error('Error:', error);
        alert('Error creating first step: ' + error.message);
    }
}

export async function deleteFirstStep(firstStepType) {
    const data = {
        firstStepType: firstStepType
    };

    try {
        const response = await fetch('/firststep/delete', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify(data)
        });

        if (!response.ok) {
            throw new Error('Network response was not ok');
        }

        const result = await response.json();
        console.log('Success:', result);
        window.location.reload(); 
    } catch (error) {
        console.error('Error deleting first step:', error);
        alert('Error deleting first step: ' + error.message);
    }
}