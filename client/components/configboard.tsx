"use client"

import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card"
import { Badge } from "@/components/ui/badge"
import { RefreshCw, Settings } from "lucide-react"

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

interface ConfigDashboardProps {
  config: Config | null
  loading: boolean
}

export default function ConfigDashboard({ config, loading }: ConfigDashboardProps) {
  return (
    <Card>
      <CardHeader>
        <CardTitle className="flex items-center space-x-2">
          <Settings className="h-5 w-5" />
          <span>Current Configuration</span>
        </CardTitle>
      </CardHeader>
      <CardContent>
        {loading ? (
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
  )
}