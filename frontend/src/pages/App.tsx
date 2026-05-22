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

type DisplayArtifact = Artifact & {
  formattedContent: string;
  extension: string;
  lineCount: number;
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

  const displayArtifacts = useMemo(() => artifacts.map(toDisplayArtifact), [artifacts]);

  const selectedArtifact = useMemo(
    () => displayArtifacts.find((artifact) => artifact.path === selectedPath) || displayArtifacts[0],
    [displayArtifacts, selectedPath],
  );

  const artifactStats = useMemo(
    () => ({
      files: displayArtifacts.length,
      lines: displayArtifacts.reduce((total, artifact) => total + artifact.lineCount, 0),
      types: new Set(displayArtifacts.map((artifact) => artifact.extension)).size,
    }),
    [displayArtifacts],
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

      const normalized = normalizeArtifacts(res.data.artifacts);
      setArtifacts(normalized);
      setSelectedPath(normalized[0]?.path || '');
      setStatus(`Generated ${normalized.length} artifact${normalized.length === 1 ? '' : 's'}`);
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
          <p className="eyebrow">Git-Ops Workspace</p>
          <h1>Pipeline Generator</h1>
          <p className="subtitle">Configure a service and review generated CI/CD files in a formatted artifact browser.</p>
        </div>
        <div className="statusPanel">
          <span className={`statusDot ${error ? 'danger' : isGenerating ? 'busy' : ''}`} />
          <span>{status}</span>
        </div>
      </section>

      <section className="workspace">
        <form className="formPanel" onSubmit={(event) => event.preventDefault()}>
          <div className="panelTitle">
            <div>
              <h2>Project Inputs</h2>
              <p>These values shape the files and deployment targets.</p>
            </div>
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
          <div className="outputHeader">
            <div className="panelTitle">
              <div>
                <h2>Generated Artifacts</h2>
                <p>Browse each generated file with readable formatting.</p>
              </div>
              <span>{artifactStats.files} files</span>
            </div>

            <div className="statsGrid">
              <div>
                <strong>{artifactStats.files}</strong>
                <span>Files</span>
              </div>
              <div>
                <strong>{artifactStats.lines}</strong>
                <span>Lines</span>
              </div>
              <div>
                <strong>{artifactStats.types}</strong>
                <span>Types</span>
              </div>
            </div>
          </div>

          {displayArtifacts.length > 0 ? (
            <div className="artifactLayout">
              <nav className="artifactList" aria-label="Generated files">
                {displayArtifacts.map((artifact) => (
                  <button
                    className={artifact.path === selectedArtifact?.path ? 'active' : ''}
                    key={artifact.path}
                    type="button"
                    onClick={() => setSelectedPath(artifact.path)}
                  >
                    <span>{artifact.path}</span>
                    <small>{artifact.lineCount} lines</small>
                  </button>
                ))}
              </nav>
              <article className="artifactPreview">
                <header>
                  <div>
                    <strong>{selectedArtifact?.path}</strong>
                    <span>{selectedArtifact?.extension.toUpperCase()} file</span>
                  </div>
                  <span>{selectedArtifact?.lineCount} lines</span>
                </header>
                <pre>{selectedArtifact?.formattedContent}</pre>
              </article>
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

function normalizeArtifacts(incoming: Artifact[]) {
  if (incoming.length !== 1) {
    return incoming;
  }

  const parsed = parseArtifactList(incoming[0].content);
  return parsed.length > 0 ? parsed : incoming;
}

function parseArtifactList(value: string): Artifact[] {
  const cleaned = stripFence(value);
  try {
    const parsed = JSON.parse(cleaned) as unknown;
    if (Array.isArray(parsed)) {
      return parsed.filter(isArtifact);
    }
    if (isRecord(parsed)) {
      const files = parsed.artifacts || parsed.files;
      return Array.isArray(files) ? files.filter(isArtifact) : [];
    }
  } catch {
    return [];
  }
  return [];
}

function isArtifact(value: unknown): value is Artifact {
  return isRecord(value) && typeof value.path === 'string' && typeof value.content === 'string';
}

function isRecord(value: unknown): value is Record<string, unknown> {
  return typeof value === 'object' && value !== null;
}

function toDisplayArtifact(artifact: Artifact): DisplayArtifact {
  const extension = extensionFor(artifact.path);
  const formattedContent = formatContent(artifact.content, extension);
  return {
    ...artifact,
    extension,
    formattedContent,
    lineCount: formattedContent.length === 0 ? 0 : formattedContent.split('\n').length,
  };
}

function extensionFor(path: string) {
  const name = path.split('/').pop() || path;
  if (name.toLowerCase() === 'dockerfile') {
    return 'dockerfile';
  }
  const extension = name.includes('.') ? name.split('.').pop() : 'txt';
  return extension?.toLowerCase() || 'txt';
}

function formatContent(content: string, extension: string) {
  const cleaned = stripFence(content).trim();
  if (['json'].includes(extension)) {
    return prettyJson(cleaned);
  }
  if (looksLikeJson(cleaned)) {
    return prettyJson(cleaned);
  }
  return cleaned.replace(/\r\n/g, '\n');
}

function prettyJson(value: string) {
  try {
    return JSON.stringify(JSON.parse(value), null, 2);
  } catch {
    return value;
  }
}

function looksLikeJson(value: string) {
  return (value.startsWith('{') && value.endsWith('}')) || (value.startsWith('[') && value.endsWith(']'));
}

function stripFence(value: string) {
  const trimmed = value.trim();
  if (!trimmed.startsWith('```')) {
    return trimmed;
  }

  const withoutOpening = trimmed.replace(/^```[a-zA-Z0-9_-]*\s*/, '');
  return withoutOpening.replace(/\s*```$/, '').trim();
}
