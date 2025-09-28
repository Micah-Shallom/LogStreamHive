"use client"

import { useState, useEffect } from "react"
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card"
import { Badge } from "@/components/ui/badge"
import { RefreshCw, Settings, FileText } from "lucide-react"

interface LogEntry {
  timestamp: string
  level: string
  message: string
  source?: string
}

interface Config {
  LOG_LEVEL: string
  LOG_FORMAT: string
  LOG_INTERVAL_SECONDS: number
}

export default function LoggerDashboard() {
  const [logs, setLogs] = useState<LogEntry[]>([])
  const [config, setConfig] = useState<Config | null>(null)
  const [logsLoading, setLogsLoading] = useState(true)
  const [configLoading, setConfigLoading] = useState(true)
  const [lastRefresh, setLastRefresh] = useState<Date>(new Date())

  // Mock data for demonstration
  const mockLogs: LogEntry[] = [
    {
      timestamp: "2024-01-15T10:30:45Z",
      level: "INFO",
      message: "Application started successfully",
      source: "main.go:45",
    },
    {
      timestamp: "2024-01-15T10:30:46Z",
      level: "DEBUG",
      message: "Database connection established",
      source: "db.go:23",
    },
    {
      timestamp: "2024-01-15T10:30:47Z",
      level: "INFO",
      message: "HTTP server listening on :8080",
      source: "server.go:67",
    },
    {
      timestamp: "2024-01-15T10:30:50Z",
      level: "WARN",
      message: "High memory usage detected: 85%",
      source: "monitor.go:12",
    },
    {
      timestamp: "2024-01-15T10:30:52Z",
      level: "ERROR",
      message: "Failed to process request: timeout exceeded",
      source: "handler.go:89",
    },
    {
      timestamp: "2024-01-15T10:30:55Z",
      level: "INFO",
      message: "Background job completed successfully",
      source: "worker.go:34",
    },
  ]

  const mockConfig: Config = {
    LOG_LEVEL: "DEBUG",
    LOG_FORMAT: "json",
    LOG_INTERVAL_SECONDS: 5,
  }

  const fetchLogs = async () => {
    setLogsLoading(true)
    try {
      // Replace with actual API call: const response = await fetch('/logs')
      // For now, using mock data
      await new Promise((resolve) => setTimeout(resolve, 500)) // Simulate API delay
      setLogs(mockLogs)
    } catch (error) {
      console.error("Failed to fetch logs:", error)
    } finally {
      setLogsLoading(false)
      setLastRefresh(new Date())
    }
  }

  const fetchConfig = async () => {
    setConfigLoading(true)
    try {
      // Replace with actual API call: const response = await fetch('/config')
      // For now, using mock data
      await new Promise((resolve) => setTimeout(resolve, 300)) // Simulate API delay
      setConfig(mockConfig)
    } catch (error) {
      console.error("Failed to fetch config:", error)
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
              <h1 className="text-2xl font-bold text-gray-900">Logger Service</h1>
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
        <div className="grid grid-cols-1 lg:grid-cols-3 gap-8">
          {/* Recent Logs Section */}
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
                          <span className="text-sm font-medium text-gray-700">Log Level</span>
                          <Badge className={getLogLevelColor(config.LOG_LEVEL)}>{config.LOG_LEVEL}</Badge>
                        </div>

                        <div className="flex justify-between items-center">
                          <span className="text-sm font-medium text-gray-700">Log Format</span>
                          <span className="text-sm text-gray-900 font-mono bg-white px-2 py-1 rounded border">
                            {config.LOG_FORMAT}
                          </span>
                        </div>

                        <div className="flex justify-between items-center">
                          <span className="text-sm font-medium text-gray-700">Interval (seconds)</span>
                          <span className="text-sm text-gray-900 font-mono bg-white px-2 py-1 rounded border">
                            {config.LOG_INTERVAL_SECONDS}
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
      </main>
    </div>
  )
}
