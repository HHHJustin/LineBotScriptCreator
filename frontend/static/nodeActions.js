// nodeActions.js

let myDiagramInstance = null;

// 初始化函數，用於設置 myDiagram 的引用
export function initialize(diagram) {
    myDiagramInstance = diagram;
}

// 顯示上下文菜單
export function showContextMenu(event, nodeData) {
    // 確保上下文菜單存在
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

    // 清空之前的內容
    menu.innerHTML = '';

    // 添加 "Add" 按鈕
    const addButton = document.createElement('button');
    addButton.textContent = 'Add';
    addButton.onclick = async () => {
        // 這裡可以調用你的添加節點函數
        await createNodeForParent(nodeData);
        hideMenu();
    };
    menu.appendChild(addButton);

    // 添加 "Remove" 按鈕
    const removeButton = document.createElement('button');
    removeButton.textContent = 'Remove';
    removeButton.onclick = () => {
        // 這裡可以調用你的刪除節點函數
        removeNode(nodeData);
        hideMenu();
    };
    menu.appendChild(removeButton);

    // 計算菜單的位置
    const menuX = event.clientX;
    const menuY = event.clientY;

    menu.style.left = `${menuX}px`;
    menu.style.top = `${menuY}px`;
    menu.style.display = 'block';

    // 點擊其他地方隱藏菜單
    document.addEventListener('click', hideMenuOutside);
}

// 隱藏上下文菜單
function hideMenu() {
    const menu = document.getElementById('contextMenu');
    if (menu) {
        menu.style.display = 'none';
    }
    document.removeEventListener('click', hideMenuOutside);
}

// 點擊外部隱藏上下文菜單
function hideMenuOutside(event) {
    const menu = document.getElementById('contextMenu');
    if (menu && !menu.contains(event.target)) {
        menu.style.display = 'none';
        document.removeEventListener('click', hideMenuOutside);
    }
}

// 根據選擇的類型創建新節點
async function createNodeForParent(parentNodeData) {
    // 這裡可以使用之前的函數來添加新節點
    const typeId = prompt("請輸入新節點的類型ID"); // 這裡可以更改為你的實現方式
    await createNode(typeId, parentNodeData);
}

// 刪除節點
function removeNode(nodeData) {
    myDiagramInstance.model.remove(nodeData);
}
