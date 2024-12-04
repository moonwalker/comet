import clsx from 'clsx';
import Heading from '@theme/Heading';
import styles from './styles.module.css';

const FeatureList = [
  {
    title: 'Easy to Use',
    src: require('@site/static/img/visual_1.webp').default,
    description: (
      <>
        Designed to streamline infrastructure provisioning and management. It provides a unified interface for handling infrastructure operations with modern tooling and practices.
      </>
    ),
  },
  {
    title: 'Focus on What Matters',
    src: require('@site/static/img/visual_2.webp').default,
    description: (
      <>
        Focus on your application logic while we handle the infrastructure complexity. Define your infrastructure as code and let automation take care of the provisioning details.
      </>
    ),
  },
  {
    title: 'Powered by Open Source',
    src: require('@site/static/img/visual_3.webp').default,
    description: (
      <>
        Built with modern technologies and best practices. Our infrastructure management tools integrate seamlessly with your existing development workflow and CI/CD pipelines.
      </>
    ),
  },
];

function Feature({src, title, description}) {
  return (
    <div className={clsx('col col--4')}>
      <div className="text--center">
        <img src={src} className={styles.featureImg} role="img" />
      </div>
      <div className="text--center padding-horiz--md">
        <Heading as="h3">{title}</Heading>
        <p>{description}</p>
      </div>
    </div>
  );
}

export default function HomepageFeatures() {
  return (
    <section className={styles.features}>
      <div className="container">
        <div className="row">
          {FeatureList.map((props, idx) => (
            <Feature key={idx} {...props} />
          ))}
        </div>
      </div>
    </section>
  );
}
