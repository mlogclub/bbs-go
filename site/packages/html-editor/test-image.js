// 简单的测试脚本，在浏览器控制台中运行
// 测试 PasteImage 扩展是否正确注册

function testPasteImageExtension() {
  // 检查编辑器是否存在
  const editorElement = document.querySelector('.ProseMirror');
  if (!editorElement) {
    console.error('❌ 未找到编辑器元素');
    return;
  }
  
  console.log('✅ 找到编辑器元素:', editorElement);
  
  // 检查编辑器是否有粘贴事件监听器
  const events = getEventListeners ? getEventListeners(editorElement) : {};
  console.log('📋 编辑器事件监听器:', events);
  
  // 模拟创建一个图片文件用于测试
  function createTestImageFile() {
    // 创建一个简单的1x1像素的PNG图片数据
    const canvas = document.createElement('canvas');
    canvas.width = 1;
    canvas.height = 1;
    const ctx = canvas.getContext('2d');
    ctx.fillStyle = 'red';
    ctx.fillRect(0, 0, 1, 1);
    
    return new Promise((resolve) => {
      canvas.toBlob((blob) => {
        const file = new File([blob], 'test.png', { type: 'image/png' });
        resolve(file);
      }, 'image/png');
    });
  }
  
  return { editorElement, createTestImageFile };
}

// 在控制台中运行这个函数
console.log('🧪 运行测试: testPasteImageExtension()');
