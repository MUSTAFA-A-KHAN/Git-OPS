import { useState } from 'react';
import axios from 'axios';

const base = import.meta.env.VITE_API_URL || 'http://localhost:8080/api/v1';

export function App() {
  const [form, setForm] = useState({ language: 'Go', framework: 'Fiber', cloudProvider: 'AWS', pipelineProvider: 'github-actions' });
  const [output, setOutput] = useState('');

  async function generate() {
    const token = (await axios.post(`${base}/auth/login`)).data.token;
    const res = await axios.post(`${base}/generate/`, form, { headers: { Authorization: `Bearer ${token}` } });
    setOutput(JSON.stringify(res.data, null, 2));
  }

  return <div className="app"><h1>Git-Ops AI Pipeline Generator</h1><button onClick={generate}>Generate</button><pre>{output}</pre></div>;
}
