<template>
  <div class="slash-commands" ref="slashCommandsRef">
    <div v-if="items.length > 0" class="search-hint">
      继续输入关键词筛选菜单...
    </div>
    <div class="slash-items-container">
      <button
        v-for="(item, index) in items"
        :key="item.title"
        class="slash-item"
        :class="{ 'is-selected': item === selectedItem }"
        @click="selectItem(item)"
        @mouseenter="handleMouseEnter(index)"
      >
        <component :is="item.icon" class="item-icon" :size="18" />
        <div class="item-content">
          <div class="item-title">
            {{ item.title }}
            <span
              v-if="item.aliases && item.aliases.length > 0"
              class="item-aliases"
            >
              /{{ item.aliases[0] }}
            </span>
          </div>
          <div class="item-description">
            {{ item.description }}
          </div>
        </div>
      </button>
    </div>
    <div v-if="items.length === 0" class="no-results">没有找到匹配的命令</div>
  </div>
</template>

<script setup lang="ts">
import { ref, watch, nextTick, onMounted, onBeforeUnmount } from "vue";
import type { CommandItem } from "./types";

const props = defineProps<{
  items: CommandItem[];
  command: (item: CommandItem) => void;
  editor: any;
  range: { from: number; to: number };
  clientRect: DOMRect | null;
}>();

const selectedItem = ref<CommandItem | null>(null);
const selectedIndex = ref(0);
const slashCommandsRef = ref<HTMLElement | null>(null);

// 处理鼠标悬停事件
const handleMouseEnter = (index: number) => {
  selectedIndex.value = index;
  selectedItem.value = props.items[index];
};

// 初始化选中第一项
if (props.items.length > 0) {
  selectedItem.value = props.items[0];
}

// 监听 items 变化，当 items 为空时清除选中项
// 当 items 有值但没有选中项时，选择第一项
// 当 items 因筛选而变化时，总是选中第一项
watch(
  () => props.items,
  (newItems) => {
    if (newItems.length === 0) {
      selectedItem.value = null;
      selectedIndex.value = 0;
    } else {
      // 始终选择第一个匹配项，无论是否已有选中项
      selectedItem.value = newItems[0];
      selectedIndex.value = 0;
      // 确保在下一个渲染周期滚动到选中项
      nextTick(() => {
        scrollToSelected();
      });
    }
  },
  { immediate: true }
);

// 滚动到选中项
const scrollToSelected = () => {
  nextTick(() => {
    const container = slashCommandsRef.value;
    const selectedElement = container?.querySelector(
      ".is-selected"
    ) as HTMLElement;

    if (container && selectedElement) {
      // 获取容器和选中元素的位置信息
      const containerRect = container.getBoundingClientRect();
      const selectedRect = selectedElement.getBoundingClientRect();

      // 判断选中元素是否在容器可视区域外
      const isAbove = selectedRect.top < containerRect.top;
      const isBelow = selectedRect.bottom > containerRect.bottom;

      if (isAbove) {
        // 如果选中项在可视区域上方，滚动到使其显示在顶部
        container.scrollTop = selectedElement.offsetTop - 8; // 添加一些间距
      } else if (isBelow) {
        // 如果选中项在可视区域下方，滚动到使其显示在底部
        container.scrollTop =
          selectedElement.offsetTop +
          selectedElement.offsetHeight -
          container.clientHeight +
          8;
      }
    }
  });
};

// 监听选中项变化，滚动到选中项
watch(
  () => selectedItem.value,
  () => {
    scrollToSelected();
  }
);

const selectItem = (item: CommandItem) => {
  selectedItem.value = item;
  props.command(item);
};

// 添加Tab键支持
const onKeyDown = (event: KeyboardEvent): boolean => {
  // 如果没有菜单项，不处理键盘事件
  if (props.items.length === 0) {
    return false;
  }

  if (event.key === "ArrowUp") {
    event.preventDefault();
    selectedIndex.value =
      (selectedIndex.value - 1 + props.items.length) % props.items.length;
    selectedItem.value = props.items[selectedIndex.value];
    return true;
  }

  if (event.key === "ArrowDown" || event.key === "Tab") {
    event.preventDefault();
    selectedIndex.value = (selectedIndex.value + 1) % props.items.length;
    selectedItem.value = props.items[selectedIndex.value];
    return true;
  }

  if (event.key === "Enter" && selectedItem.value) {
    event.preventDefault();
    selectItem(selectedItem.value);
    return true;
  }

  return false;
};

// 快捷键，1-9 数字选择对应的项目
const handleNumberShortcuts = (event: KeyboardEvent) => {
  if (props.items.length === 0) return;

  const num = parseInt(event.key);
  // 检查键是否为 1-9 之间的数字
  if (!isNaN(num) && num >= 1 && num <= 9) {
    // 数字键 1-9 对应索引 0-8
    const index = num - 1;
    if (index < props.items.length) {
      event.preventDefault();
      selectedIndex.value = index;
      selectedItem.value = props.items[index];
      selectItem(props.items[index]);
    }
  }
};

// 添加全局键盘事件监听
onMounted(() => {
  document.addEventListener("keydown", handleNumberShortcuts);
});

onBeforeUnmount(() => {
  document.removeEventListener("keydown", handleNumberShortcuts);
});

defineExpose({
  onKeyDown,
});
</script>

<style lang="scss">
.slash-commands {
  background: white;
  border-radius: 8px;
  box-shadow: 0 0 0 1px rgba(0, 0, 0, 0.05), 0px 10px 20px rgba(0, 0, 0, 0.1);
  padding: 0.5rem;
  width: 280px;
  max-height: 400px;
  overflow-y: auto;
  animation: popup-fade 0.15s ease-out;

  // 美化滚动条样式
  &::-webkit-scrollbar {
    width: 4px;
    opacity: 0;
    transition: opacity 0.3s;
  }

  &:hover::-webkit-scrollbar {
    opacity: 1;
  }

  &::-webkit-scrollbar-track {
    background: transparent;
    border-radius: 3px;
  }

  &::-webkit-scrollbar-thumb {
    background-color: rgba(0, 0, 0, 0.2);
    border-radius: 3px;
    transition: background-color 0.3s;
    min-height: 30px;
    max-height: 60px;

    &:hover {
      background-color: rgba(0, 0, 0, 0.3);
    }
  }

  // Firefox滚动条样式
  scrollbar-width: thin;
  scrollbar-color: transparent transparent;

  &:hover {
    scrollbar-color: rgba(0, 0, 0, 0.2) transparent;
  }

  /* 暗色主题滚动条 */
  .dark & {
    scrollbar-color: transparent transparent;

    &:hover {
      scrollbar-color: rgba(255, 255, 255, 0.2) transparent;
    }

    &::-webkit-scrollbar-thumb {
      background-color: rgba(255, 255, 255, 0);

      &:hover {
        background-color: rgba(255, 255, 255, 0.3);
      }
    }

    &:hover::-webkit-scrollbar-thumb {
      background-color: rgba(255, 255, 255, 0.2);
    }
  }
}

.search-hint {
  font-size: 12px;
  color: #6b7280;
  margin-bottom: 0.5rem;
  padding: 0 0.5rem;
}

.slash-items-container {
  max-height: 320px;
  overflow-y: auto;
}

.slash-item {
  display: flex;
  align-items: center;
  padding: 0.625rem;
  margin: 0.125rem 0;
  width: 100%;
  background: none;
  border: none;
  border-radius: 6px;
  cursor: pointer;
  text-align: left;
  transition: all 0.15s ease;

  &:hover {
    background-color: #f5f7fa;
  }

  &.is-selected {
    background-color: #f0f2f5;

    .item-icon {
      color: #4b5563;
    }

    .item-title {
      color: #000;
    }
  }

  .item-icon {
    margin-right: 0.75rem;
    color: #6b7280;
    flex-shrink: 0;
    transition: color 0.15s ease;
  }

  .item-content {
    flex: 1;
    overflow: hidden;
  }

  .item-title {
    font-size: 14px;
    color: #111827;
    display: flex;
    justify-content: space-between;
    align-items: center;
    transition: color 0.15s ease;

    .item-aliases {
      font-size: 12px;
      color: #6b7280;
      margin-left: 4px;
      border-radius: 4px;
      padding: 1px 4px;
      background-color: #f3f4f6;
    }
  }

  .item-description {
    font-size: 12px;
    color: #6b7280;
    white-space: nowrap;
    overflow: hidden;
    text-overflow: ellipsis;
  }
}

.no-results {
  padding: 1rem;
  text-align: center;
  color: #6b7280;
  font-size: 14px;
}

/* 动画 */
@keyframes popup-fade {
  from {
    opacity: 0;
    transform: translateY(-5px);
  }
  to {
    opacity: 1;
    transform: translateY(0);
  }
}

/* 暗色主题样式 */
.dark {
  .slash-commands {
    background: #1f2937;
    box-shadow: 0 0 0 1px rgba(255, 255, 255, 0.05),
      0px 10px 20px rgba(0, 0, 0, 0.2);
  }

  .no-results {
    color: #9ca3af;
  }

  .slash-item {
    &:hover {
      background-color: #2d3748;
    }

    &.is-selected {
      background-color: #374151;

      .item-icon {
        color: #e5e7eb;
      }

      .item-title {
        color: #f9fafb;
      }
    }

    .item-icon {
      color: #9ca3af;
    }

    .item-title {
      color: #f9fafb;

      .item-aliases {
        background-color: #374151;
        color: #9ca3af;
      }
    }

    .item-description {
      color: #9ca3af;
    }
  }
}
</style>
