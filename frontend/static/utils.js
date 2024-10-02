export function getCurrentIDAndTypeFromURL() {
    const url = window.location.pathname;  
    const segments = url.split('/');  
    const currentID = segments[3];    

    return {
        currentID: currentID,
    };
}