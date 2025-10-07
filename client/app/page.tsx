"use client"

import { useState, useEffect } from "react"
import { Centrifuge } from "centrifuge";
import { RefreshCw, FileText, Settings } from "lucide-react"
import { Tabs, TabsContent, TabsList, TabsTrigger } from "@/components/ui/tabs"
import { Button } from "@/components/ui/button"
import {
  Dialog,
  DialogContent,
  DialogHeader,
  DialogTitle,
  DialogTrigger,
} from "@/components/ui/dialog"
import Statistics from "@/components/statisticsboard"
import ConfigDashboard from "@/components/configboard"
import LoggerDashboard from "@/components/logsboard"
import CollectorDashboard from "@/components/collectorboard"

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
  const [collectorLogs, setCollectorLogs] = useState<LogEntry[]>([])
  const [config, setConfig] = useState<Config | null>(null)
  const [stats, setStats] = useState<Statistics | null>(null)
  
  const [logsLoading, setLogsLoading] = useState(true)
  const [configLoading, setConfigLoading] = useState(true)
  const [statsLoading, setStatsLoading] = useState(true)
  
  const [lastRefresh, setLastRefresh] = useState<Date>(new Date())

  const API_URL = process.env.NEXT_PUBLIC_API_URL || "http://localhost:8080"
  const WS_URL = process.env.NEXT_PUBLIC_WS_URL || "ws://centrifugo:8000/connection/websocket"

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
    
    const centrifuge = new Centrifuge("ws://localhost:8000/connection/websocket", {
    token: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM3MjIiLCJleHAiOjE2NTU0NDgyOTl9.mUU9s5kj3yqp-SAEqloGy8QBgsLg0llA7lKUNwtHRnw"
  });

    centrifuge.on('connecting', function (ctx) {
      console.log(`connecting: ${ctx.code}, ${ctx.reason}`);
    }).on('connected', function (ctx) {
      console.log(`connected over ${ctx.transport}`);
    }).on('disconnected', function (ctx) {
      console.log(`disconnected: ${ctx.code}, ${ctx.reason}`);
    }).connect();

    const sub = centrifuge.newSubscription("logs");

    sub.on('publication', function (ctx) {
      setCollectorLogs((prevLogs) => [ctx.data, ...prevLogs])
    }).on('subscribing', function (ctx) {
      console.log(`subscribing: ${ctx.code}, ${ctx.reason}`);
    }).on('subscribed', function (ctx) {
      console.log('subscribed', ctx);
    }).on('unsubscribed', function (ctx) {
      console.log(`unsubscribed: ${ctx.code}, ${ctx.reason}`);
    }).subscribe();

    return () => {
      clearInterval(interval)
      centrifuge.disconnect()
    }
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
              <Dialog>
                <DialogTrigger asChild>
                  <Button variant="outline" size="sm">
                    <Settings className="h-4 w-4 mr-2" />
                    View Config
                  </Button>
                </DialogTrigger>
                <DialogContent className="max-w-4xl">
                  <DialogHeader>
                    <DialogTitle>Current Configuration</DialogTitle>
                  </DialogHeader>
                  <ConfigDashboard 
                    config={config} 
                    loading={configLoading}
                  />
                </DialogContent>
              </Dialog>
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
            <div className="grid grid-cols-1 lg:grid-cols-2 gap-8 mt-4">
              <div className="lg:col-span-1">
                <LoggerDashboard 
                  logs={logs} 
                  loading={logsLoading}
                  onRefresh={fetchLogs}
                />
              </div>
              <div className="lg:col-span-1">
                <CollectorDashboard 
                  logs={collectorLogs}
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