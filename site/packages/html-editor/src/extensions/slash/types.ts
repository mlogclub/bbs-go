import type { Component } from 'vue'
import type { Editor } from '@tiptap/core'
import {
    Heading1,
    Heading2,
    List,
    ListOrdered,
    Quote,
    Code2,
    Image as ImageIcon,
    MinusSquare,
    Table as TableIcon,
    Link as LinkIcon,
} from 'lucide-vue-next'

export interface CommandItem {
    title: string
    description: string
    icon: Component
    aliases?: string[]
    command: ({ editor, range }: { editor: Editor; range: { from: number; to: number } }) => void
}

export const getSuggestionItems = () => [
    {
        title: '标题1',
        description: '大标题',
        icon: Heading1,
        aliases: ['h1', '一级标题', 'heading1'],
        command: ({ editor, range }: { editor: Editor; range: { from: number; to: number } }) => {
            editor.chain().focus().deleteRange(range).setNode('heading', { level: 1 }).run()
        },
    },
    {
        title: '标题2',
        description: '二级标题',
        icon: Heading2,
        aliases: ['h2', '二级标题', 'heading2'],
        command: ({ editor, range }: { editor: Editor; range: { from: number; to: number } }) => {
            editor.chain().focus().deleteRange(range).setNode('heading', { level: 2 }).run()
        },
    },
    {
        title: '标题3',
        description: '三级标题',
        icon: Heading2,
        aliases: ['h3', '三级标题', 'heading3'],
        command: ({ editor, range }: { editor: Editor; range: { from: number; to: number } }) => {
            editor.chain().focus().deleteRange(range).setNode('heading', { level: 3 }).run()
        },
    },
    {
        title: '无序列表',
        description: '创建无序列表',
        icon: List,
        aliases: ['ul', 'bullet', 'list'],
        command: ({ editor, range }: { editor: Editor; range: { from: number; to: number } }) => {
            editor.chain().focus().deleteRange(range).toggleBulletList().run()
        },
    },
    {
        title: '有序列表',
        description: '创建有序列表',
        icon: ListOrdered,
        aliases: ['ol', 'ordered', 'numbered'],
        command: ({ editor, range }: { editor: Editor; range: { from: number; to: number } }) => {
            editor.chain().focus().deleteRange(range).toggleOrderedList().run()
        },
    },
    {
        title: '引用',
        description: '插入引用文本',
        icon: Quote,
        aliases: ['quote', 'blockquote', '引用文本'],
        command: ({ editor, range }: { editor: Editor; range: { from: number; to: number } }) => {
            editor.chain().focus().deleteRange(range).toggleBlockquote().run()
        },
    },
    {
        title: '代码块',
        description: '插入代码块',
        icon: Code2,
        aliases: ['code', 'codeblock', '代码'],
        command: ({ editor, range }: { editor: Editor; range: { from: number; to: number } }) => {
            editor.chain().focus().deleteRange(range).toggleCodeBlock().run()
        },
    },
    // TODO 
    // {
    //     title: '图片',
    //     description: '插入图片',
    //     icon: ImageIcon,
    //     aliases: ['img', 'image', 'picture'],
    //     command: ({ editor, range }: { editor: Editor; range: { from: number; to: number } }) => {
    //         const url = window.prompt('输入图片URL')
    //         if (url) {
    //             editor.chain().focus().deleteRange(range).setResizableImage({ src: url }).run()
    //         }
    //     },
    // },
    // TODO
    // {
    //     title: '链接',
    //     description: '插入链接',
    //     icon: LinkIcon,
    //     aliases: ['link', 'url', 'href'],
    //     command: ({ editor, range }: { editor: Editor; range: { from: number; to: number } }) => {
    //         editor.chain().focus().deleteRange(range).openLinkDialog().run()
    //     },
    // },
    {
        title: '分割线',
        description: '插入水平分割线',
        icon: MinusSquare,
        aliases: ['hr', 'line', 'divider'],
        command: ({ editor, range }: { editor: Editor; range: { from: number; to: number } }) => {
            editor.chain().focus().deleteRange(range).setHorizontalRule().run()
        },
    },
    // TODO
    // {
    //     title: '表格',
    //     description: '插入表格',
    //     icon: TableIcon,
    //     aliases: ['table', 'grid', 'tb'],
    //     command: ({ editor, range }: { editor: Editor; range: { from: number; to: number } }) => {
    //         editor
    //             .chain()
    //             .focus()
    //             .deleteRange(range)
    //             .insertTable({ rows: 3, cols: 3, withHeaderRow: true })
    //             .run()
    //     },
    // },
]