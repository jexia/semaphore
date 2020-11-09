module.exports = {
  docs: [
    {
      type: 'category',
      label: 'Installation',
      items: [
        'installation',
        'installation.cli',
        'installation.docker',
        'installation.source',
        'installation.contributing',
      ],
    },
    {
      type: 'category',
      label: 'Flows',
      items: [
        'flows',
        'flows.references',
        'flows.conditions',
        'flows.requests',
        'flows.errors',
        'flows.proxy',
        'flows.rollbacks',
      ],
    },
    {
      type: 'doc',
      id: 'functions',
    },
    {
      type: 'doc',
      id: 'devops',
      label: 'Infrastructure',
      items: [
        "service_discovery.configuration",
      ]
    }
  ]
};
