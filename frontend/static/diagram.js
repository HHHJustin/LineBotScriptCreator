export function createDiagram($) {
    if (window.myDiagram) {
        return window.myDiagram;
      }
    
      const myDiagram = $(go.Diagram, 'myDiagramDiv');
    
      // Diagram 的其他配置...
    
      window.myDiagram = myDiagram;
      return myDiagram;
}
