import React from 'react';
import clsx from 'clsx';
import Layout from '@theme/Layout';
import Link from '@docusaurus/Link';
import useDocusaurusContext from '@docusaurus/useDocusaurusContext';
import useBaseUrl from '@docusaurus/useBaseUrl';
import styles from './styles.module.css';

const features = [
  {
    title: 'Connect with anything',
    imageUrl: 'img/undraw_code_typing_7jnv.svg',
    description: (
      <>
        Use the right tool for the job.
        Semaphore supports various protocols out of the box with the ability to supporting additional protocols through modules.
        Endpoints could be created to expose a single flow through multiple protocols.
      </>
    ),
  },
  {
    title: 'Blazing fast',
    imageUrl: 'img/undraw_Outer_space_drqu.svg',
    description: (
      <>
        Semaphore scales up to your needs.
        Branches are created to execute resources concurrently.
        Branches are based on dependencies between resources made through references or hard coded values.
        Creating high-performance flows is almost boringly easy.
      </>
    ),
  },
  {
    title: 'Adapts to your environment',
    imageUrl: 'img/undraw_connected_world_wuay.svg',
    description: (
      <>
        Semaphore integrates with your existing system(s).
        Define flows through simple and strict typed definitions.
        Use your already existing schema definitions such as Protobuffers.
        Or extend Semaphore with custom modules and proprietary software.
        Integrate services through flow definitions and create a great experience for your customers and your teams.
      </>
    ),
  },
];

function Feature({imageUrl, title, description}) {
  const imgUrl = useBaseUrl(imageUrl);
  return (
    <div className={clsx('col col--4', styles.feature)}>
      {imgUrl && (
        <div className="text--center">
          <img className={styles.featureImage} src={imgUrl} alt={title} />
        </div>
      )}
      <h3>{title}</h3>
      <p>{description}</p>
    </div>
  );
}

function Home() {
  const context = useDocusaurusContext();
  const {siteConfig = {}} = context;
  return (
    <Layout
      title={`Hello from ${siteConfig.title}`}
      description="Description will go into a meta tag in <head />">
      <header className={clsx('hero hero--primary', styles.heroBanner)}>
        <div className="container">
          <h1 className="hero__title">{siteConfig.title}</h1>
          <p className="hero__subtitle">{siteConfig.tagline}</p>
          <div className={styles.buttons}>
            <Link
              className={clsx(
                'button button--outline button--secondary button--lg',
                styles.getStarted,
              )}
              to={useBaseUrl('docs/installation')}>
              Get Started
            </Link>
          </div>
        </div>
      </header>
      <main>
        {features && features.length > 0 && (
          <section className={styles.features}>
            <div className="container">
              <div className="row">
                {features.map((props, idx) => (
                  <Feature key={idx} {...props} />
                ))}
              </div>
            </div>
          </section>
        )}
      </main>
    </Layout>
  );
}

export default Home;
