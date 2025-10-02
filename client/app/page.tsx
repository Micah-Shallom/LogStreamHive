"use client"

import { useState, useEffect } from "react"
import { RefreshCw, FileText } from "lucide-react"
import { Tabs, TabsContent, TabsList, TabsTrigger } from "@/components/ui/tabs"
import Statistics from "@/components/statisticsboard"
import ConfigDashboard from "@/components/configboard"
import LoggerDashboard from "@/components/logsboard"

interface LogEntry {
  timestamp: string
  log_type: string
  user_id: string
  duration: number
  message: string
  request_id: string
  service: string
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

interface Statistics {
  logTypeCounts: Record<string, number>
  serviceDurations: Record<string, number>
  serviceCallCounts: Record<string, number>
  errorSequences: ErrorSequence[]
  anomalyDetections: Anomaly[]
  updatedAt: string
}

interface ErrorSequence {
  startTime: string
  endTime: string
  count: number
  service: string
}

interface Anomaly {
  timestamp: string
  service: string
  metricName: string
  value: number
  threshold: number
}

export default function Dashboard() {
  const [logs, setLogs] = useState<LogEntry[]>([])
  const [config, setConfig] = useState<Config | null>(null)
  const [stats, setStats] = useState<Statistics | null>(null)
  
  const [logsLoading, setLogsLoading] = useState(true)
  const [configLoading, setConfigLoading] = useState(true)
  const [statsLoading, setStatsLoading] = useState(true)
  
  const [lastRefresh, setLastRefresh] = useState<Date>(new Date())

  const API_URL = process.env.NEXT_PUBLIC_API_URL || "http://localhost:8000"

  const fetchLogs = async () => {
    setLogsLoading(true)
    try {
      const response = await fetch(`${API_URL}/logs`)
      if (!response.ok) {
        console.error("Response not ok:", response.statusText)
        throw new Error(`HTTP error! status: ${response.status}`)
      }
      const res = await response.json()
      if (res.status !== "success") {
        throw new Error(res.Message || "Unknown error")
      }

      const data: LogEntry[] = res.data || []
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

      const res = await response.json()
      if (res.status !== "success") {
        throw new Error(res.Message || "Unknown error")
      }
      
      const data: Config = res.data 
      setConfig(data)
    } catch (error) {
      console.error("Failed to fetch config:", error)
      setConfig(null)
    } finally {
      setConfigLoading(false)
    }
  }

  const fetchStatistics = async () => {
    setStatsLoading(true)
    try {
      const response = await fetch(`${API_URL}/statistics`)
      if (!response.ok) {
        throw new Error(`HTTP error! status: ${response.status}`)
      }

      const res = await response.json()
      if (res.status !== "success") {
        throw new Error(res.Message || "Unknown error")
      }

      const data: Statistics = res.data

      if (!data.logTypeCounts || !data.serviceDurations || !data.serviceCallCounts) {
        throw new Error("Incomplete statistics data received from server.")
      }

      setStats(data)
    } catch (error) {
      console.error("Failed to fetch statistics:", error)
      setStats(null)
    } finally {
      setStatsLoading(false)
    }
  }

  const refreshAll = () => {
    fetchLogs()
    fetchConfig()
    fetchStatistics()
  }

  useEffect(() => {
    // Initial fetch of all data
    fetchLogs()
    fetchConfig()
    fetchStatistics()

    // Auto-refresh logs every 10 seconds
    const interval = setInterval(fetchLogs, 10000)
    return () => clearInterval(interval)
  }, [])

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
                onClick={refreshAll}
                disabled={logsLoading || configLoading || statsLoading}
                className="inline-flex items-center px-3 py-2 border border-gray-300 shadow-sm text-sm leading-4 font-medium rounded-md text-gray-700 bg-white hover:bg-gray-50 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-blue-500 disabled:opacity-50"
              >
                <RefreshCw className={`h-4 w-4 mr-2 ${(logsLoading || configLoading || statsLoading) ? "animate-spin" : ""}`} />
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
                <LoggerDashboard 
                  logs={logs} 
                  loading={logsLoading}
                  onRefresh={fetchLogs}
                />
              </div>
              <div className="lg:col-span-1">
                <ConfigDashboard 
                  config={config} 
                  loading={configLoading}
                />
              </div>
            </div>
          </TabsContent>
          
          <TabsContent value="statistics">
            <Statistics 
              stats={stats}
              loading={statsLoading}
              onRefresh={fetchStatistics}
            />
          </TabsContent>
        </Tabs>
      </main>
    </div>
  )
}