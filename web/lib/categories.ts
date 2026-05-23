import type { Category } from "@/lib/api/types"

export function flattenCategories(categories: Category[] = []): Category[] {
  const result: Category[] = []

  function walk(list: Category[]) {
    list.forEach((category) => {
      result.push(category)
      if (category.children?.length) {
        walk(category.children)
      }
    })
  }

  walk(Array.isArray(categories) ? categories : [])
  return result
}

export function filterCategoryTree(
  categories: Category[] = [],
  predicate: (item: Category) => boolean
): Category[] {
  if (!Array.isArray(categories) || !categories.length) {
    return []
  }

  return categories.reduce<Category[]>((result, category) => {
    const children = filterCategoryTree(category.children || [], predicate)
    if (!predicate(category) && !children.length) {
      return result
    }
    result.push({ ...category, children })
    return result
  }, [])
}

export function hasCategory(categories: Category[] = [], categoryId: number | string) {
  const targetId = Number(categoryId)
  return flattenCategories(categories).some((node) => Number(node.id) === targetId)
}

export function getFirstCategoryId(categories: Category[] = []) {
  const first = flattenCategories(categories)[0]
  return first ? Number(first.id) : 0
}
