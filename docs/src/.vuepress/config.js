module.exports = {
  title: 'Semaphore',
  description: 'Build powerfull data flows',
  base: '/semaphore/',
  themeConfig: {
    lastUpdated: false,
    sidebar: [
      "/Guide/",
      "/Get_Started/",
      "/Cookbook/",
      "/DevOps/"
    ],
    /* nav: [{
      text: "Guide",
      link: "/guide/"
    }], */
    repo: "jexia/semaphore",
    repoLabel: "Contribute",
    displayAllHeaders: false
  },
  plugins: ['@vuepress/medium-zoom']
}