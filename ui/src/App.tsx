import { useState, useRef, useEffect } from 'react'
import ReactMarkdown from 'react-markdown'
import { Send, Upload, FileText, Settings, Bot, User, BrainCircuit, Sparkles, Loader2, Database } from 'lucide-react'

interface Message {
  role: 'user' | 'assistant'
  content: string
}

export default function App() {
  const [messages, setMessages] = useState<Message[]>([])
  const [input, setInput] = useState('')
  const [isStreaming, setIsStreaming] = useState(false)
  const [uploading, setUploading] = useState(false)
  const messagesEndRef = useRef<HTMLDivElement>(null)

  const scrollToBottom = () => {
    messagesEndRef.current?.scrollIntoView({ behavior: 'smooth' })
  }

  useEffect(() => {
    scrollToBottom()
  }, [messages])

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault()
    if (!input.trim() || isStreaming) return

    const userMsg = input
    setInput('')
    setMessages(prev => [...prev, { role: 'user', content: userMsg }])
    setMessages(prev => [...prev, { role: 'assistant', content: '' }])
    setIsStreaming(true)

    try {
      const res = await fetch(`http://localhost:8080/api/stream?q=${encodeURIComponent(userMsg)}`)
      if (!res.body) throw new Error('No body')
      
      const reader = res.body.getReader()
      const decoder = new TextDecoder()
      
      while (true) {
        const { done, value } = await reader.read()
        if (done) break
        
        const chunk = decoder.decode(value)
        const lines = chunk.split('\n')
        
        for (const line of lines) {
          if (line.startsWith('data: ')) {
            const data = line.slice(6)
            setMessages(prev => {
              const newMsgs = [...prev]
              newMsgs[newMsgs.length - 1].content += data
              return newMsgs
            })
          }
        }
      }
    } catch (err) {
      console.error(err)
      setMessages(prev => {
        const newMsgs = [...prev]
        newMsgs[newMsgs.length - 1].content = "Error connecting to CortexRAG engine."
        return newMsgs
      })
    } finally {
      setIsStreaming(false)
    }
  }

  const handleFileUpload = async (e: React.ChangeEvent<HTMLInputElement>) => {
    const file = e.target.files?.[0]
    if (!file) return

    setUploading(true)
    const formData = new FormData()
    formData.append('file', file)

    try {
      const res = await fetch('http://localhost:8080/api/ingest', {
        method: 'POST',
        body: formData,
      })
      if (!res.ok) throw new Error("Upload failed")
      
      // Simulate slight delay for effect
      await new Promise(r => setTimeout(r, 800))
      
      alert('Document ingrained into vector store successfully!')
    } catch (err) {
      console.error('Upload failed', err)
      alert('Failed to ingest document. Ensure the backend is running.')
    } finally {
      setUploading(false)
    }
  }

  return (
    <div className="flex h-screen bg-[var(--background)] text-[var(--foreground)] font-sans overflow-hidden">
      
      {/* Sidebar - Glassmorphism */}
      <div className="w-72 bg-[var(--sidebar)]/80 backdrop-blur-xl border-r border-[var(--glass-border)] flex flex-col relative z-10 shadow-2xl">
        <div className="p-6 border-b border-[var(--glass-border)] flex items-center gap-3">
          <div className="w-10 h-10 rounded-xl bg-gradient-to-br from-blue-500 to-purple-600 flex items-center justify-center shadow-lg shadow-blue-500/20">
            <BrainCircuit className="w-5 h-5 text-white animate-pulse-slow" />
          </div>
          <div>
            <h1 className="font-bold text-lg tracking-tight text-white">CortexRAG</h1>
            <p className="text-xs text-[var(--muted-foreground)] font-medium">Vector Engine Gateway</p>
          </div>
        </div>
        
        <div className="flex-1 p-4 overflow-y-auto">
          <div className="flex items-center justify-between mb-4 px-2">
            <h2 className="text-xs font-bold text-[var(--muted-foreground)] uppercase tracking-widest flex items-center gap-2">
              <Database className="w-3.5 h-3.5" />
              Knowledge Base
            </h2>
          </div>
          
          <div className="space-y-1">
            <div className="group flex items-center gap-3 text-sm p-3 rounded-lg hover:bg-white/5 transition-all cursor-pointer border border-transparent hover:border-[var(--glass-border)]">
              <div className="w-8 h-8 rounded-lg bg-[var(--muted)] flex items-center justify-center group-hover:bg-blue-500/10 transition-colors">
                <FileText className="w-4 h-4 text-[var(--muted-foreground)] group-hover:text-blue-400 transition-colors" />
              </div>
              <div className="flex-1 overflow-hidden">
                <p className="truncate font-medium text-slate-300 group-hover:text-white transition-colors">No documents ingested</p>
                <p className="text-xs text-slate-500">Upload to begin</p>
              </div>
            </div>
          </div>
        </div>

        <div className="p-4 border-t border-[var(--glass-border)] bg-gradient-to-t from-[var(--sidebar)] to-transparent">
          <label className="group relative flex items-center justify-center gap-2 w-full p-3 bg-gradient-to-r from-blue-600 to-indigo-600 hover:from-blue-500 hover:to-indigo-500 rounded-xl cursor-pointer transition-all shadow-lg shadow-blue-900/20 hover:shadow-blue-900/40 border border-blue-400/20 hover:-translate-y-0.5">
            {uploading ? (
              <Loader2 className="w-5 h-5 text-white animate-spin" />
            ) : (
              <Upload className="w-5 h-5 text-white group-hover:scale-110 transition-transform" />
            )}
            <span className="font-semibold text-sm text-white">
              {uploading ? 'Ingesting...' : 'Ingest Document'}
            </span>
            <input type="file" className="hidden" accept=".md,.pdf,.txt" onChange={handleFileUpload} disabled={uploading} />
          </label>
        </div>
      </div>

      {/* Main Content */}
      <div className="flex-1 flex flex-col relative">
        {/* Header */}
        <header className="h-16 border-b border-[var(--glass-border)] flex items-center justify-between px-8 bg-[var(--background)]/40 backdrop-blur-md sticky top-0 z-20">
          <div className="flex items-center gap-3">
            <Sparkles className="w-4 h-4 text-purple-400" />
            <span className="text-sm font-medium text-slate-300">Active Model:</span>
            <select className="bg-[var(--muted)]/50 backdrop-blur-sm border border-[var(--glass-border)] rounded-lg px-3 py-1.5 text-sm font-medium text-white outline-none focus:border-blue-500 focus:ring-1 focus:ring-blue-500 transition-all cursor-pointer hover:bg-[var(--muted)]">
              <option value="gemini">Gemini 1.5 Flash (via RPC)</option>
              <option value="ollama">Ollama (Local Inference)</option>
            </select>
          </div>
          <button className="p-2 hover:bg-white/10 rounded-xl transition-all border border-transparent hover:border-[var(--glass-border)]">
            <Settings className="w-5 h-5 text-slate-400 hover:text-white transition-colors" />
          </button>
        </header>

        {/* Chat Area */}
        <main className="flex-1 overflow-y-auto p-8 space-y-8 scroll-smooth">
          {messages.length === 0 ? (
            <div className="h-full flex flex-col items-center justify-center text-center max-w-2xl mx-auto animate-fade-in">
              <div className="w-24 h-24 rounded-3xl bg-gradient-to-br from-blue-500/10 to-purple-600/10 border border-[var(--glass-border)] flex items-center justify-center mb-8 shadow-2xl relative">
                <div className="absolute inset-0 bg-blue-500 blur-3xl opacity-10 rounded-full" />
                <BrainCircuit className="w-12 h-12 text-blue-400 relative z-10" />
              </div>
              <h2 className="text-4xl font-bold mb-4 bg-clip-text text-transparent bg-gradient-to-r from-white to-slate-400 tracking-tight">
                High-Throughput RAG
              </h2>
              <p className="text-slate-400 text-lg leading-relaxed max-w-lg mx-auto font-medium">
                Upload your PDFs or Markdown in the sidebar. CortexRAG will chunk, embed, and index them via HNSW for ultra-fast semantic retrieval.
              </p>
            </div>
          ) : (
            messages.map((msg, i) => (
              <div key={i} className={`flex gap-4 animate-fade-in ${msg.role === 'user' ? 'justify-end' : 'justify-start'}`}>
                {msg.role === 'assistant' && (
                  <div className="w-10 h-10 rounded-xl bg-gradient-to-br from-slate-800 to-slate-900 border border-[var(--glass-border)] flex items-center justify-center shrink-0 shadow-lg mt-1">
                    <Bot className="w-5 h-5 text-blue-400" />
                  </div>
                )}
                
                <div className={`max-w-[75%] rounded-2xl p-5 shadow-xl ${
                  msg.role === 'user' 
                    ? 'bg-blue-600 text-white rounded-tr-sm border border-blue-500/50' 
                    : 'bg-[#0f172a]/80 backdrop-blur-md border border-[var(--glass-border)] rounded-tl-sm text-slate-200'
                }`}>
                  {msg.role === 'user' ? (
                    <p className="whitespace-pre-wrap leading-relaxed">{msg.content}</p>
                  ) : (
                    <div className="prose prose-invert max-w-none prose-p:leading-relaxed prose-pre:bg-black/50 prose-pre:border prose-pre:border-slate-800 prose-headings:font-semibold">
                      {msg.content === '' ? (
                        <div className="flex items-center gap-2 text-blue-400 font-medium">
                          <Loader2 className="w-4 h-4 animate-spin" />
                          <span>Retrieving context & generating...</span>
                        </div>
                      ) : (
                        <ReactMarkdown>{msg.content}</ReactMarkdown>
                      )}
                    </div>
                  )}
                </div>

                {msg.role === 'user' && (
                  <div className="w-10 h-10 rounded-xl bg-slate-800 border border-slate-700 flex items-center justify-center shrink-0 shadow-lg mt-1">
                    <User className="w-5 h-5 text-slate-300" />
                  </div>
                )}
              </div>
            ))
          )}
          <div ref={messagesEndRef} className="h-4" />
        </main>

        {/* Input Area */}
        <div className="p-6 bg-gradient-to-t from-[var(--background)] via-[var(--background)]/90 to-transparent pt-12 relative z-20">
          <form onSubmit={handleSubmit} className="max-w-4xl mx-auto relative group">
            <div className="absolute -inset-1 bg-gradient-to-r from-blue-500 to-purple-600 rounded-2xl blur opacity-20 group-hover:opacity-30 transition duration-1000 group-hover:duration-200"></div>
            <div className="relative flex items-center bg-[#0a0f1c] border border-[var(--glass-border)] rounded-2xl shadow-2xl">
              <input
                type="text"
                value={input}
                onChange={e => setInput(e.target.value)}
                placeholder="Ask a question across your vector index..."
                className="w-full bg-transparent text-white placeholder-slate-500 rounded-2xl pl-6 pr-16 py-5 focus:outline-none focus:ring-1 focus:ring-blue-500/50 transition-all font-medium text-lg"
                disabled={isStreaming}
              />
              <button
                type="submit"
                disabled={!input.trim() || isStreaming}
                className="absolute right-3 p-3 bg-blue-600 hover:bg-blue-500 disabled:bg-slate-800 disabled:text-slate-500 disabled:cursor-not-allowed rounded-xl transition-all shadow-lg text-white group/btn"
              >
                {isStreaming ? (
                  <Loader2 className="w-5 h-5 animate-spin" />
                ) : (
                  <Send className="w-5 h-5 group-hover/btn:-translate-y-0.5 group-hover/btn:translate-x-0.5 transition-transform" />
                )}
              </button>
            </div>
          </form>
          <p className="text-center text-xs text-slate-500 font-medium mt-4 flex items-center justify-center gap-1.5">
            <Sparkles className="w-3.5 h-3.5" />
            CortexRAG engine processes documents locally via configured vector backends.
          </p>
        </div>
      </div>
    </div>
  )
}
