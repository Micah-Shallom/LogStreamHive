"use client"

import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card"
import { Badge } from "@/components/ui/badge"
import { RefreshCw, FileText } from "lucide-react"

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

interface LoggerDashboardProps {
  logs: LogEntry[]
  loading: boolean
  onRefresh: () => void
}

export default function LoggerDashboard({ logs, loading, onRefresh }: LoggerDashboardProps) {
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
    <Card>
      <CardHeader>
        <CardTitle className="flex items-center space-x-2">
          <FileText className="h-5 w-5" />
          <span>Recent Logs</span>
        </CardTitle>
      </CardHeader>
      <CardContent>
        <div className="h-96 overflow-y-auto space-y-3 bg-gray-50 p-4 rounded-lg">
          {loading ? (
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
                      <Badge className={getLogLevelColor(log.log_type)}>{log.log_type}</Badge>
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
  )
}