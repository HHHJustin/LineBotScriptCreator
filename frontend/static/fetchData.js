export async function fetchData(myDiagram) {
    try {
        const response = await fetch('/nodes/get');
        if (!response.ok) {
            throw new Error(`HTTP error! status: ${response.status}`);
        }
        const data = await response.json();
        myDiagram.model = new go.GraphLinksModel(data.nodes, data.links);
    } catch (error) {
        console.error('抓取資料失敗:', error);
    }
}
