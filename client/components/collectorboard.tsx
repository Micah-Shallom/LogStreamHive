"use client"

import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card"
import { List } from "lucide-react"

interface ParsedLog {
  timestamp: string
  log_type: string
  user_id: string
  duration: number
  message: string
  request_id: string
  service: string
}

export interface RawLog {
  timestamp: string
  file_path: string
  line: string
}

interface ProcessedLog {
  id: number
  isError?: boolean
  line: ParsedLog | { message: string }
}

interface CollectorDashboardProps {
  logs: RawLog[]
}

// Utility: determine text color based on log type
const getLogTypeColor = (logType?: string): string => {
  if (!logType) return "text-gray-500"
  switch (logType.toUpperCase()) {
    case "INFO":
      return "text-blue-500"
    case "WARNING":
      return "text-yellow-500"
    case "ERROR":
      return "text-red-500"
    case "CRITICAL":
      return "text-red-700 font-bold"
    default:
      return "text-gray-500"
  }
}

export default function CollectorDashboard({ logs }: CollectorDashboardProps) {
  const parsedLogs: ProcessedLog[] = logs.map((log, index) => {
    try {
      const parsedLine: ParsedLog = JSON.parse(log.line)
      return {
        id: index,
        line: parsedLine,
      }
    } catch (error) {
      console.error("Failed to parse log line:", log.line, error)
      return {
        id: index,
        line: {
          message: `Unparseable log line: ${log.line}`,
        },
        isError: true,
      }
    }
  })

  return (
    <Card>
      <CardHeader>
        <CardTitle className="flex items-center space-x-2">
          <List className="h-5 w-5" />
          <span>Collector Board</span>
        </CardTitle>
      </CardHeader>
      <CardContent>
        <div className="space-y-4">
          {parsedLogs.length > 0 ? (
            parsedLogs.slice(0, 10).map((log) => {
              const line = log.line as ParsedLog
              return (
                <div
                  key={log.id}
                  className="bg-gray-100 p-3 rounded-lg shadow-sm font-mono text-xs"
                >
                  {!log.isError && "log_type" in line ? (
                    <>
                      <div className="flex items-center space-x-2 mb-1">
                        <span className={`font-bold ${getLogTypeColor(line.log_type)}`}>
                          [{line.log_type}]
                        </span>
                        <span className="text-gray-500">
                          {new Date(line.timestamp).toLocaleString()}
                        </span>
                        <span className="font-semibold text-purple-600">{line.service}</span>
                      </div>
                      <p className="text-gray-800">{line.message}</p>
                      <div className="flex flex-wrap text-gray-500 mt-1">
                        <span>UserID: {line.user_id}</span>
                        <span className="mx-2">|</span>
                        <span>Duration: {line.duration}ms</span>
                        <span className="mx-2">|</span>
                        <span>ReqID: {line.request_id}</span>
                      </div>
                    </>
                  ) : (
                    <div>
                      <span className="font-bold text-red-500">[PARSE ERROR]</span>
                      <pre className="text-xs whitespace-pre-wrap">
                        {(log.line as any).message}
                      </pre>
                    </div>
                  )}
                </div>
              )
            })
          ) : (
            <div className="text-center py-8 text-gray-500">
              No logs received yet.
            </div>
          )}
        </div>
      </CardContent>
    </Card>
  )
}
