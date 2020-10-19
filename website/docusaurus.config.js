module.exports = {
  title: 'Semaphore',
  tagline: 'A straightforward micro-service orchestrator',
  url: 'https://jexia.github.io',
  baseUrl: '/semaphore/',
  onBrokenLinks: 'throw',
  favicon: 'img/favicon.ico',
  organizationName: 'jexia', // Usually your GitHub org/user name.
  projectName: 'semaphore', // Usually your repo name.
  themeConfig: {
    navbar: {
      title: 'Semaphore',
      logo: {
        alt: 'Semaphore',
        src: 'img/logo.svg',
      },
      items: [
        {
          to: 'docs/installation/',
          activeBasePath: 'docs',
          label: 'Docs',
          position: 'left',
        },
        // {to: 'blog', label: 'Blog', position: 'left'},
        {
          href: 'https://github.com/jexia/semaphore',
          label: 'GitHub',
          position: 'right',
        },
      ],
    },
    footer: {
      style: 'dark',
      links: [
        {
          title: 'Docs',
          items: [
            {
              label: 'Installation',
              to: 'docs/installation/',
            },
          ],
        },
        {
          title: 'Community',
          items: [
            {
              label: 'Discord',
              href: 'https://discord.com/invite/qWByeWG',
            },
          ],
        },
        {
          title: 'More',
          items: [
            // {
            //   label: 'Blog',
            //   to: 'blog',
            // },
            {
              label: 'GitHub',
              href: 'https://github.com/jexia/semaphore',
            },
          ],
        },
      ],
    },
  },
  presets: [
    [
      '@docusaurus/preset-classic',
      {
        docs: {
          sidebarPath: require.resolve('./sidebars.js'),
          // Please change this to your repo.
          editUrl:
            'https://github.com/jexia/semaphore/edit/master/website/',
        },
        blog: {
          showReadingTime: true,
          // Please change this to your repo.
          editUrl:
            'https://github.com/jexia/semaphore/edit/master/website/blog/',
        },
        theme: {
          customCss: require.resolve('./src/css/custom.css'),
        },
      },
    ],
  ],
};
