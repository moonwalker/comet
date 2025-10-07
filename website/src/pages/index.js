import clsx from 'clsx';
import Link from '@docusaurus/Link';
import useDocusaurusContext from '@docusaurus/useDocusaurusContext';
import Layout from '@theme/Layout';
import Heading from '@theme/Heading';
import styles from './index.module.css';

function HomepageHeader() {
  const {siteConfig} = useDocusaurusContext();
  return (
    <header className={clsx('hero', styles.heroBanner)}>
      <div className="container">
        <div className={styles.heroContent}>
          <div className={styles.heroText}>
            <span className={styles.badge}>🚀 Simple. Powerful. Open.</span>
            <Heading as="h1" className={styles.heroTitle}>
              Infrastructure management made <span className={styles.highlight}>simple</span>
            </Heading>
            <p className={styles.heroSubtitle}>
              A cosmic tool for provisioning and managing infrastructure.
              Define your infra with code, automatic backend generation, and built-in secrets management.
            </p>
            <div className={styles.buttons}>
              <Link
                className="button button--primary button--lg"
                to="/docs/intro">
                Get Started →
              </Link>
              <Link
                className="button button--secondary button--lg"
                to="/docs/comparison">
                Why Comet?
              </Link>
            </div>
            <div className={styles.quickStart}>
              <code>curl -fsSL https://moonwalker.github.io/comet/install.sh | sh</code>
            </div>
          </div>
        </div>
      </div>
    </header>
  );
}

function Feature({icon, title, description}) {
  return (
    <div className={clsx('col col--4')}>
      <div className={styles.feature}>
        <div className={styles.featureIcon}>{icon}</div>
        <Heading as="h3">{title}</Heading>
        <p>{description}</p>
      </div>
    </div>
  );
}

function FeaturesSection() {
  const features = [
    {
      icon: '🚀',
      title: 'JavaScript Configuration',
      description: 'Define infrastructure using familiar JavaScript instead of limited HCL. Leverage the full power of a programming language.',
    },
    {
      icon: '🔄',
      title: 'Auto Backend Generation',
      description: 'Comet generates backend.tf.json files automatically. No more manual backend configuration.',
    },
    {
      icon: '🔗',
      title: 'Cross-Stack References',
      description: 'Simple state() function to reference outputs from other stacks. Dependencies made easy.',
    },
    {
      icon: '🔐',
      title: 'Built-in Secrets',
      description: 'Native SOPS integration for encrypted secrets. Manage sensitive data securely out of the box.',
    },
    {
      icon: '📦',
      title: 'Component Reusability',
      description: 'Define components once, reuse across environments with different configurations.',
    },
    {
      icon: '⚡',
      title: 'OpenTofu & Terraform',
      description: 'Works seamlessly with both OpenTofu and Terraform. Minimal abstraction, maximum compatibility.',
    },
  ];

  return (
    <section className={styles.features}>
      <div className="container">
        <div className="row">
          {features.map((props, idx) => (
            <Feature key={idx} {...props} />
          ))}
        </div>
      </div>
    </section>
  );
}

function CodeExample() {
  return (
    <section className={styles.codeSection}>
      <div className="container">
        <div className="row">
          <div className="col col--6">
            <Heading as="h2">Write infrastructure in JavaScript</Heading>
            <p className={styles.codeDescription}>
              Define your infrastructure using JavaScript's expressive syntax.
              Share configuration, use variables, and leverage the ecosystem you already know.
            </p>
            <ul className={styles.benefitsList}>
              <li>✅ Familiar syntax and tooling</li>
              <li>✅ Powerful templating with Go templates</li>
              <li>✅ DRY configuration across environments</li>
              <li>✅ Type hints with JSDoc</li>
            </ul>
          </div>
          <div className="col col--6">
            <div className={styles.codeBlock}>
              <pre>{`// stacks/production.js
const { settings } = require('./shared.js')

stack('production', { settings })

backend('gcs', {
  bucket: 'terraform-state',
  prefix: 'comet/{{ .stack }}/{{ .component }}'
})

const vpc = component('vpc', 'modules/vpc', {
  cidr_block: '10.0.0.0/16',
  region: '{{ .settings.region }}'
})

const gke = component('gke', 'modules/gke', {
  network: vpc.id,
  cluster_name: 'prod-cluster'
})`}</pre>
            </div>
          </div>
        </div>
      </div>
    </section>
  );
}

function ComparisonSection() {
  return (
    <section className={styles.comparison}>
      <div className="container">
        <div className="text--center margin-bottom--lg">
          <Heading as="h2">How does Comet compare?</Heading>
          <p className={styles.comparisonSubtitle}>
            Comet fills the gap between plain Terraform and heavy enterprise frameworks
          </p>
        </div>
        <div className={styles.comparisonGrid}>
          <div className={styles.comparisonCard}>
            <h3>Plain OpenTofu/Terraform</h3>
            <p className={styles.useCase}>Best for: Simple setups</p>
            <ul>
              <li>❌ Manual backend config</li>
              <li>❌ Verbose multi-env setup</li>
              <li>❌ Manual cross-stack refs</li>
              <li>✅ No abstractions</li>
            </ul>
          </div>
          <div className={clsx(styles.comparisonCard, styles.highlight)}>
            <div className={styles.recommendedBadge}>Recommended for most teams</div>
            <h3>Comet</h3>
            <p className={styles.useCase}>Best for: Small-medium teams</p>
            <ul>
              <li>✅ JavaScript config</li>
              <li>✅ Auto backend generation</li>
              <li>✅ Built-in SOPS secrets</li>
              <li>✅ Minimal abstraction</li>
            </ul>
          </div>
          <div className={styles.comparisonCard}>
            <h3>Terragrunt / Atmos</h3>
            <p className={styles.useCase}>Best for: Large enterprises</p>
            <ul>
              <li>✅ Battle-tested</li>
              <li>✅ Large community</li>
              <li>⚠️ More complex</li>
              <li>⚠️ Steeper learning curve</li>
            </ul>
          </div>
        </div>
        <div className="text--center margin-top--lg">
          <Link className="button button--outline button--primary" to="/docs/comparison">
            See Detailed Comparison →
          </Link>
        </div>
      </div>
    </section>
  );
}

function CTASection() {
  return (
    <section className={styles.cta}>
      <div className="container">
        <div className="text--center">
          <Heading as="h2">Ready to simplify your infrastructure?</Heading>
          <p className={styles.ctaSubtitle}>
            Get started with Comet in minutes. Build from source or check out the documentation.
          </p>
          <div className={styles.buttons}>
            <Link
              className="button button--primary button--lg"
              to="/docs/intro">
              Get Started
            </Link>
            <Link
              className="button button--secondary button--lg"
              href="https://github.com/moonwalker/comet">
              View on GitHub ⭐
            </Link>
          </div>
        </div>
      </div>
    </section>
  );
}

export default function Home() {
  const {siteConfig} = useDocusaurusContext();
  return (
    <Layout
      title="Home"
      description="Cosmic tool for provisioning and managing infrastructure with JavaScript-based configuration">
      <HomepageHeader />
      <main>
        <FeaturesSection />
        <CodeExample />
        <ComparisonSection />
        <CTASection />
      </main>
    </Layout>
  );
}
