"use client"

import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card"
import { List } from "lucide-react"

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

interface CollectorDashboardProps {
  logs: LogEntry[]
}

export default function CollectorDashboard({ logs }: CollectorDashboardProps) {
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
          {logs.length > 0 ? (
            logs.map((log, index) => (
              <div key={index} className="bg-gray-50 p-2 rounded-lg">
                <pre className="text-xs">{JSON.stringify(log, null, 2)}</pre>
              </div>
            ))
          ) : (
            <div className="text-center py-8 text-gray-500">No logs received yet.</div>
          )}
        </div>
      </CardContent>
    </Card>
  )
}