import { useMemo, useState } from 'react';
import axios from 'axios';

const base = import.meta.env.VITE_API_URL || 'http://localhost:8080/api/v1';

type GenerationForm = {
  projectName: string;
  language: string;
  framework: string;
  cloudProvider: string;
  deploymentTarget: string;
  repositoryType: string;
  testingFramework: string;
  buildTool: string;
  containerizationPreference: string;
  kubernetesUsage: boolean;
  monitoringRequirements: string;
  pipelineProvider: string;
};

type Artifact = {
  path: string;
  content: string;
};

const initialForm: GenerationForm = {
  projectName: 'orders-service',
  language: 'Go',
  framework: 'Fiber',
  cloudProvider: 'AWS',
  deploymentTarget: 'EKS',
  repositoryType: 'monorepo',
  testingFramework: 'Go test',
  buildTool: 'Go modules',
  containerizationPreference: 'Docker',
  kubernetesUsage: true,
  monitoringRequirements: 'Prometheus metrics and structured logs',
  pipelineProvider: 'github-actions',
};

const fieldOptions: Partial<Record<keyof GenerationForm, string[]>> = {
  language: ['Go', 'Node.js', 'Python', 'Java', '.NET'],
  framework: ['Fiber', 'Express', 'FastAPI', 'Spring Boot', 'ASP.NET Core'],
  cloudProvider: ['AWS', 'Azure', 'GCP', 'On-prem'],
  deploymentTarget: ['EKS', 'AKS', 'GKE', 'Kubernetes', 'VM'],
  repositoryType: ['monorepo', 'single service', 'polyrepo'],
  testingFramework: ['Go test', 'Jest', 'Pytest', 'JUnit', 'xUnit'],
  buildTool: ['Go modules', 'npm', 'Poetry', 'Maven', 'Gradle'],
  containerizationPreference: ['Docker', 'Buildpacks', 'None'],
  pipelineProvider: ['github-actions', 'gitlab-ci', 'jenkins', 'azure-devops'],
};

export function App() {
  const [form, setForm] = useState<GenerationForm>(initialForm);
  const [artifacts, setArtifacts] = useState<Artifact[]>([]);
  const [selectedPath, setSelectedPath] = useState('');
  const [status, setStatus] = useState('Ready');
  const [error, setError] = useState('');
  const [isGenerating, setIsGenerating] = useState(false);

  const selectedArtifact = useMemo(
    () => artifacts.find((artifact) => artifact.path === selectedPath) || artifacts[0],
    [artifacts, selectedPath],
  );

  function updateField<K extends keyof GenerationForm>(key: K, value: GenerationForm[K]) {
    setForm((current) => ({ ...current, [key]: value }));
  }

  async function generate() {
    setError('');
    setStatus('Authenticating');
    setIsGenerating(true);

    try {
      const token = (await axios.post(`${base}/auth/login`)).data.token;
      setStatus('Generating artifacts');

      const res = await axios.post<{ artifacts: Artifact[] }>(`${base}/generate/`, form, {
        headers: { Authorization: `Bearer ${token}` },
      });

      setArtifacts(res.data.artifacts);
      setSelectedPath(res.data.artifacts[0]?.path || '');
      setStatus(`Generated ${res.data.artifacts.length} artifact${res.data.artifacts.length === 1 ? '' : 's'}`);
    } catch (err) {
      const message = axios.isAxiosError(err)
        ? err.response?.data?.message || err.response?.data || err.message
        : 'Generation failed';
      setError(String(message));
      setStatus('Generation failed');
    } finally {
      setIsGenerating(false);
    }
  }

  return (
    <main className="shell">
      <section className="header">
        <div>
          <p className="eyebrow">Git-Ops AI</p>
          <h1>Pipeline Generator</h1>
          <p className="subtitle">Generate CI/CD pipelines, Dockerfiles, and Kubernetes manifests for a service.</p>
        </div>
        <div className="statusPanel">
          <span className={`statusDot ${error ? 'danger' : isGenerating ? 'busy' : ''}`} />
          <span>{status}</span>
        </div>
      </section>

      <section className="workspace">
        <form className="formPanel" onSubmit={(event) => event.preventDefault()}>
          <div className="panelTitle">
            <h2>Project Inputs</h2>
            <button type="button" onClick={generate} disabled={isGenerating}>
              {isGenerating ? 'Generating...' : 'Generate'}
            </button>
          </div>

          <label>
            Project name
            <input value={form.projectName} onChange={(event) => updateField('projectName', event.target.value)} />
          </label>

          <div className="grid">
            {Object.entries(fieldOptions).map(([key, options]) => (
              <label key={key}>
                {labelFor(key)}
                <select
                  value={String(form[key as keyof GenerationForm])}
                  onChange={(event) => updateField(key as keyof GenerationForm, event.target.value as never)}
                >
                  {options.map((option) => (
                    <option key={option}>{option}</option>
                  ))}
                </select>
              </label>
            ))}
          </div>

          <label className="checkRow">
            <input
              checked={form.kubernetesUsage}
              type="checkbox"
              onChange={(event) => updateField('kubernetesUsage', event.target.checked)}
            />
            Include Kubernetes manifests
          </label>

          <label>
            Monitoring requirements
            <textarea
              value={form.monitoringRequirements}
              onChange={(event) => updateField('monitoringRequirements', event.target.value)}
            />
          </label>

          {error && <div className="errorBox">{error}</div>}
        </form>

        <section className="outputPanel">
          <div className="panelTitle">
            <h2>Generated Artifacts</h2>
            <span>{artifacts.length} files</span>
          </div>

          {artifacts.length > 0 ? (
            <div className="artifactLayout">
              <nav className="artifactList" aria-label="Generated files">
                {artifacts.map((artifact) => (
                  <button
                    className={artifact.path === selectedArtifact?.path ? 'active' : ''}
                    key={artifact.path}
                    type="button"
                    onClick={() => setSelectedPath(artifact.path)}
                  >
                    {artifact.path}
                  </button>
                ))}
              </nav>
              <pre>{selectedArtifact?.content}</pre>
            </div>
          ) : (
            <div className="emptyState">
              <strong>No artifacts yet</strong>
              <span>Fill in the project inputs and run Generate.</span>
            </div>
          )}
        </section>
      </section>
    </main>
  );
}

function labelFor(key: string) {
  return key.replace(/([A-Z])/g, ' $1').replace(/^./, (char) => char.toUpperCase());
}
