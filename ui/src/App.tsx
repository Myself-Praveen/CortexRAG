import { useState, useRef, useEffect } from 'react'
import ReactMarkdown from 'react-markdown'
import { Send, Upload, FileText, Settings, Bot, User, BrainCircuit } from 'lucide-react'

interface Message {
  role: 'user' | 'assistant'
  content: string
}

export default function App() {
  const [messages, setMessages] = useState<Message[]>([])
  const [input, setInput] = useState('')
  const [isStreaming, setIsStreaming] = useState(false)
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
    
    // Add empty assistant message to stream into
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
    } finally {
      setIsStreaming(false)
    }
  }

  const handleFileUpload = async (e: React.ChangeEvent<HTMLInputElement>) => {
    const file = e.target.files?.[0]
    if (!file) return

    const formData = new FormData()
    formData.append('file', file)

    try {
      await fetch('http://localhost:8080/api/ingest', {
        method: 'POST',
        body: formData,
      })
      alert('Document ingested successfully!')
    } catch (err) {
      console.error('Upload failed', err)
      alert('Upload failed')
    }
  }

  return (
    <div className="flex h-screen bg-[#09090b] text-[#fafafa] font-sans">
      {/* Sidebar */}
      <div className="w-64 border-r border-[#27272a] bg-[#18181b] flex flex-col">
        <div className="p-4 border-b border-[#27272a] flex items-center gap-2">
          <BrainCircuit className="w-6 h-6 text-blue-500" />
          <h1 className="font-bold text-lg">CortexRAG</h1>
        </div>
        
        <div className="flex-1 p-4 overflow-y-auto">
          <h2 className="text-xs font-semibold text-[#a1a1aa] uppercase tracking-wider mb-4">Documents</h2>
          <div className="space-y-2">
            {/* Document list placeholder */}
            <div className="flex items-center gap-2 text-sm p-2 rounded hover:bg-[#27272a] transition-colors cursor-pointer text-[#a1a1aa]">
              <FileText className="w-4 h-4" />
              <span>No documents yet</span>
            </div>
          </div>
        </div>

        <div className="p-4 border-t border-[#27272a]">
          <label className="flex items-center justify-center gap-2 w-full p-2 bg-[#27272a] hover:bg-[#3f3f46] rounded-md cursor-pointer transition-colors text-sm">
            <Upload className="w-4 h-4" />
            <span>Upload Document</span>
            <input type="file" className="hidden" accept=".md,.pdf" onChange={handleFileUpload} />
          </label>
        </div>
      </div>

      {/* Main Content */}
      <div className="flex-1 flex flex-col">
        {/* Header */}
        <header className="h-14 border-b border-[#27272a] flex items-center justify-between px-6 bg-[#09090b]/80 backdrop-blur-sm">
          <div className="flex items-center gap-2">
            <span className="text-sm text-[#a1a1aa]">Model:</span>
            <select className="bg-transparent border border-[#27272a] rounded px-2 py-1 text-sm outline-none focus:border-blue-500">
              <option value="gemini">Gemini 1.5 Flash</option>
              <option value="ollama">Ollama (Llama 3)</option>
            </select>
          </div>
          <button className="p-2 hover:bg-[#27272a] rounded-full transition-colors">
            <Settings className="w-5 h-5 text-[#a1a1aa]" />
          </button>
        </header>

        {/* Chat Area */}
        <main className="flex-1 overflow-y-auto p-6 space-y-6">
          {messages.length === 0 ? (
            <div className="h-full flex flex-col items-center justify-center text-center opacity-50">
              <BrainCircuit className="w-16 h-16 mb-4 text-[#a1a1aa]" />
              <h2 className="text-2xl font-bold mb-2">Welcome to CortexRAG</h2>
              <p className="max-w-md">Upload documents in the sidebar, then ask questions about them. The engine will retrieve relevant context and stream the answer.</p>
            </div>
          ) : (
            messages.map((msg, i) => (
              <div key={i} className={`flex gap-4 ${msg.role === 'user' ? 'justify-end' : 'justify-start'}`}>
                {msg.role === 'assistant' && (
                  <div className="w-8 h-8 rounded-full bg-blue-500/20 flex items-center justify-center shrink-0">
                    <Bot className="w-5 h-5 text-blue-500" />
                  </div>
                )}
                
                <div className={`max-w-[80%] rounded-2xl p-4 ${
                  msg.role === 'user' 
                    ? 'bg-blue-600 text-white rounded-tr-sm' 
                    : 'bg-[#18181b] border border-[#27272a] rounded-tl-sm'
                }`}>
                  {msg.role === 'user' ? (
                    <p className="whitespace-pre-wrap">{msg.content}</p>
                  ) : (
                    <div className="prose prose-invert prose-sm max-w-none">
                      {msg.content === '' ? (
                        <span className="animate-pulse">Thinking...</span>
                      ) : (
                        <ReactMarkdown>{msg.content}</ReactMarkdown>
                      )}
                    </div>
                  )}
                </div>

                {msg.role === 'user' && (
                  <div className="w-8 h-8 rounded-full bg-[#27272a] flex items-center justify-center shrink-0">
                    <User className="w-5 h-5 text-[#a1a1aa]" />
                  </div>
                )}
              </div>
            ))
          )}
          <div ref={messagesEndRef} />
        </main>

        {/* Input Area */}
        <div className="p-4 border-t border-[#27272a] bg-[#09090b]">
          <form onSubmit={handleSubmit} className="max-w-4xl mx-auto relative flex items-center">
            <input
              type="text"
              value={input}
              onChange={e => setInput(e.target.value)}
              placeholder="Ask anything about your documents..."
              className="w-full bg-[#18181b] border border-[#27272a] rounded-full pl-6 pr-12 py-4 focus:outline-none focus:border-blue-500 transition-colors"
              disabled={isStreaming}
            />
            <button
              type="submit"
              disabled={!input.trim() || isStreaming}
              className="absolute right-2 p-2 bg-blue-600 hover:bg-blue-700 disabled:bg-blue-600/50 disabled:cursor-not-allowed rounded-full transition-colors"
            >
              <Send className="w-5 h-5 text-white" />
            </button>
          </form>
          <p className="text-center text-xs text-[#a1a1aa] mt-3">
            CortexRAG processes documents locally or via Gemini API. Verification is recommended.
          </p>
        </div>
      </div>
    </div>
  )
}
