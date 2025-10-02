"use client"

import { useState, useEffect } from "react"
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card"
import { Badge } from "@/components/ui/badge"
import { RefreshCw, Settings, FileText } from "lucide-react"
import { Tabs, TabsContent, TabsList, TabsTrigger } from "@/components/ui/tabs"
import Statistics from "@/components/statistics"

interface LogEntry {
  timestamp: string
  log_type: string
  user_id: string
  duration: number
  message: string
  request_id: string
  service: string
  level: string
  source?: string
}

interface Config {
  LOG_RATE: number
  LOG_TYPES: string[]
  LOG_DISTRIBUTION: Record<string, number>
  OUTPUT_FILE: string
  CONSOLE_OUTPUT: boolean
  LOG_FORMAT: string
  SERVICES: string[]
  ENABLE_BURSTS: boolean
  BURST_FREQUENCY: number
  BURST_MULTIPLIER: number
  BURST_DURATION: number
}

export default function LoggerDashboard() {
  const [logs, setLogs] = useState<LogEntry[]>([])
  const [config, setConfig] = useState<Config | null>(null)
  const [logsLoading, setLogsLoading] = useState(true)
  const [configLoading, setConfigLoading] = useState(true)
  const [lastRefresh, setLastRefresh] = useState<Date>(new Date())

  const API_URL = process.env.NEXT_PUBLIC_API_URL || "http://localhost:8000"

  const fetchLogs = async () => {
    setLogsLoading(true)
    try {
      const response = await fetch(`${API_URL}/logs`)
      if (!response.ok){
        console.error("Response not ok:", response.statusText)
        throw new Error(`HTTP error! status: ${response.status}`)
      }
      const data: LogEntry[] = await response.json()
      console.log("Fetched logs:", data)
      setLogs(data.reverse())
    } catch (error) {
      console.error("Failed to fetch logs:", error)
      setLogs([])
    } finally {
      setLogsLoading(false)
      setLastRefresh(new Date())
    }
  }

  const fetchConfig = async () => {
    setConfigLoading(true)
    try {
      const response = await fetch(`${API_URL}/config`)
      if (!response.ok) {
        console.error("Response not ok:", response.statusText)
        throw new Error(`HTTP error! status: ${response.status}`)
      }

      const data: Config = await response.json()
      console.log("Fetched config:", data)
      setConfig(data)
    } catch (error) {
      console.error("Failed to fetch config:", error)
      setConfig(null)
    } finally {
      setConfigLoading(false)
    }
  }

  useEffect(() => {
    fetchLogs()
    fetchConfig()

    // Auto-refresh logs every 10 seconds
    const interval = setInterval(fetchLogs, 10000)
    return () => clearInterval(interval)
  }, [])

  const getLogLevelColor = (level: string) => {
    switch (level.toLowerCase()) {
      case "error":
        return "bg-red-100 text-red-800 border-red-200"
      case "warn":
        return "bg-yellow-100 text-yellow-800 border-yellow-200"
      case "info":
        return "bg-blue-100 text-blue-800 border-blue-200"
      case "debug":
        return "bg-gray-100 text-gray-800 border-gray-200"
      default:
        return "bg-gray-100 text-gray-800 border-gray-200"
    }
  }

  const getLogTypeColor = (type: string) => {
    switch (type.toLowerCase()) {
      case "error":
        return "bg-red-100 text-red-800 border-red-200"
      case "warn":
      case "warning":
        return "bg-yellow-100 text-yellow-800 border-yellow-200"
      case "info":
        return "bg-blue-100 text-blue-800 border-blue-200"
      case "debug":
        return "bg-gray-100 text-gray-800 border-gray-200"
      default:
        return "bg-gray-100 text-gray-800 border-gray-200"
    }
  }

  const formatTimestamp = (timestamp: string) => {
    return new Date(timestamp).toLocaleString()
  }

  return (
    <div className="min-h-screen bg-gray-50">
      {/* Header */}
      <header className="bg-white shadow-sm border-b">
        <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
          <div className="flex justify-between items-center h-16">
            <div className="flex items-center space-x-3">
              <FileText className="h-8 w-8 text-blue-600" />
              <h1 className="text-2xl font-bold text-gray-900">LogStreamHive</h1>
            </div>
            <div className="flex items-center space-x-4">
              <span className="text-sm text-gray-500">Last updated: {lastRefresh.toLocaleTimeString()}</span>
              <button
                onClick={fetchLogs}
                disabled={logsLoading}
                className="inline-flex items-center px-3 py-2 border border-gray-300 shadow-sm text-sm leading-4 font-medium rounded-md text-gray-700 bg-white hover:bg-gray-50 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-blue-500 disabled:opacity-50"
              >
                <RefreshCw className={`h-4 w-4 mr-2 ${logsLoading ? "animate-spin" : ""}`} />
                Refresh
              </button>
            </div>
          </div>
        </div>
      </header>

      <main className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
        <Tabs defaultValue="logs" className="w-full">
          <TabsList className="grid w-full grid-cols-2">
            <TabsTrigger value="logs">Logs</TabsTrigger>
            <TabsTrigger value="statistics">Statistics</TabsTrigger>
          </TabsList>
          
          <TabsContent value="logs">
            <div className="grid grid-cols-1 lg:grid-cols-3 gap-8 mt-4">
              <div className="lg:col-span-2">
                <Card>
                  <CardHeader>
                    <CardTitle className="flex items-center space-x-2">
                      <FileText className="h-5 w-5" />
                      <span>Recent Logs</span>
                    </CardTitle>
                  </CardHeader>
                  <CardContent>
                    <div className="h-96 overflow-y-auto space-y-3 bg-gray-50 p-4 rounded-lg">
                      {logsLoading ? (
                        <div className="flex items-center justify-center h-full">
                          <RefreshCw className="h-6 w-6 animate-spin text-gray-400" />
                          <span className="ml-2 text-gray-500">Loading logs...</span>
                        </div>
                      ) : logs.length === 0 ? (
                        <div className="flex items-center justify-center h-full text-gray-500">No logs available</div>
                      ) : (
                        logs.map((log, index) => (
                          <div
                            key={index}
                            className="bg-white p-3 rounded-md border border-gray-200 hover:shadow-sm transition-shadow"
                          >
                            <div className="flex items-start justify-between">
                              <div className="flex-1">
                                <div className="flex items-center space-x-2 mb-1">
                                  <Badge className={getLogLevelColor(log.level)}>{log.level}</Badge>
                                  <span className="text-xs text-gray-500">{formatTimestamp(log.timestamp)}</span>
                                  {log.source && <span className="text-xs text-gray-400">{log.source}</span>}
                                </div>
                                <p className="text-sm text-gray-900 font-mono">{log.message}</p>
                              </div>
                            </div>
                          </div>
                        ))
                      )}
                    </div>
                  </CardContent>
                </Card>
              </div>

              {/* Current Configuration Section */}
              <div className="lg:col-span-1">
                <Card>
                  <CardHeader>
                    <CardTitle className="flex items-center space-x-2">
                      <Settings className="h-5 w-5" />
                      <span>Current Configuration</span>
                    </CardTitle>
                  </CardHeader>
                  <CardContent>
                    {configLoading ? (
                      <div className="flex items-center justify-center py-8">
                        <RefreshCw className="h-6 w-6 animate-spin text-gray-400" />
                        <span className="ml-2 text-gray-500">Loading config...</span>
                      </div>
                    ) : config ? (
                      <div className="space-y-4">
                        <div className="bg-gray-50 p-4 rounded-lg">
                          <div className="space-y-3">
                            <div className="flex justify-between items-center">
                              <span className="text-sm font-medium text-gray-700">Log Rate</span>
                              <span className="text-sm text-gray-900 font-mono bg-white px-2 py-1 rounded border">
                                {config.LOG_RATE} logs/sec
                              </span>
                            </div>

                            <div className="flex justify-between items-center">
                              <span className="text-sm font-medium text-gray-700">Log Types</span>
                              <span className="text-sm text-gray-900 font-mono bg-white px-2 py-1 rounded border">
                                {config.LOG_TYPES.join(", ")}
                              </span>
                            </div>

                            <div className="flex flex-col">
                              <span className="text-sm font-medium text-gray-700 mb-1">Log Distribution</span>
                              <div className="space-y-1">
                                {Object.entries(config.LOG_DISTRIBUTION).map(([level, weight]) => (
                                  <div key={level} className="flex justify-between items-center">
                                    <span className="text-xs font-medium text-gray-600">{level}</span>
                                    <span className="text-xs font-mono bg-white px-2 py-0.5 rounded border">
                                      {weight}
                                    </span>
                                  </div>
                                ))}
                              </div>
                            </div>

                            <div className="flex justify-between items-center">
                              <span className="text-sm font-medium text-gray-700">Output File</span>
                              <span className="text-sm text-gray-900 font-mono bg-white px-2 py-1 rounded border truncate max-w-[200px]">
                                {config.OUTPUT_FILE}
                              </span>
                            </div>

                            <div className="flex justify-between items-center">
                              <span className="text-sm font-medium text-gray-700">Console Output</span>
                              <Badge variant={config.CONSOLE_OUTPUT ? "default" : "secondary"}>
                                {config.CONSOLE_OUTPUT ? "Enabled" : "Disabled"}
                              </Badge>
                            </div>

                            <div className="flex justify-between items-center">
                          <span className="text-sm font-medium text-gray-700">Log Format</span>
                          <span className="text-sm text-gray-900 font-mono bg-white px-2 py-1 rounded border">
                            {config.LOG_FORMAT}
                          </span>
                        </div>

                        <div className="flex justify-between items-center">
                          <span className="text-sm font-medium text-gray-700">Services</span>
                          <span className="text-sm text-gray-900 font-mono bg-white px-2 py-1 rounded border">
                            {config.SERVICES.join(", ")}
                          </span>
                        </div>

                        <div className="flex justify-between items-center">
                          <span className="text-sm font-medium text-gray-700">Enable Bursts</span>
                          <Badge variant={config.ENABLE_BURSTS ? "default" : "secondary"}>
                            {config.ENABLE_BURSTS ? "Enabled" : "Disabled"}
                          </Badge>
                        </div>

                        <div className="flex justify-between items-center">
                          <span className="text-sm font-medium text-gray-700">Burst Frequency</span>
                          <span className="text-sm text-gray-900 font-mono bg-white px-2 py-1 rounded border">
                            {config.BURST_FREQUENCY}
                          </span>
                        </div>

                        <div className="flex justify-between items-center">
                          <span className="text-sm font-medium text-gray-700">Burst Multiplier</span>
                          <span className="text-sm text-gray-900 font-mono bg-white px-2 py-1 rounded border">
                            {config.BURST_MULTIPLIER}
                          </span>
                        </div>

                        <div className="flex justify-between items-center">
                          <span className="text-sm font-medium text-gray-700">Burst Duration</span>
                          <span className="text-sm text-gray-900 font-mono bg-white px-2 py-1 rounded border">
                            {config.BURST_DURATION}s
                          </span>
                        </div>
                          </div>
                        </div>

                        <div className="text-xs text-gray-500 mt-4">
                          Configuration is read from environment variables and updated on service restart.
                        </div>
                      </div>
                    ) : (
                      <div className="text-center py-8 text-gray-500">Failed to load configuration</div>
                    )}
                  </CardContent>
                </Card>
              </div>
            </div>
          </TabsContent>
          
          <TabsContent value="statistics">
            <Statistics />
          </TabsContent>
        </Tabs>
      </main>
    </div>
  )
}