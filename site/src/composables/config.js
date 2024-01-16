export function useSiteTitle(...subTitles) {
  const configStore = useConfigStore();
  const titleParts = [];
  if (subTitles && subTitles.length) {
    titleParts.push(...subTitles);
  }
  if (configStore.config.siteTitle) {
    titleParts.push(configStore.config.siteTitle);
  }
  return titleParts.join(" - ");
}

export function useSiteDescription() {
  const configStore = useConfigStore();
  return configStore.config.siteDescription;
}

export function useSiteKeywords() {
  const configStore = useConfigStore();
  return configStore.config.siteKeywords;
}

export function useTopicSiteTitle(topic) {
  if (topic.type === 0) {
    return useSiteTitle(topic.title);
  } else {
    return useSiteTitle(topic.content);
  }
}
