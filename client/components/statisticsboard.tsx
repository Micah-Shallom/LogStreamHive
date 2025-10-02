"use client"

import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card"
import { BarChart, Bar, XAxis, YAxis, CartesianGrid, Tooltip, Legend, ResponsiveContainer } from 'recharts'
import { RefreshCw } from "lucide-react"

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

interface Statistics {
  logTypeCounts: Record<string, number>
  serviceDurations: Record<string, number>
  serviceCallCounts: Record<string, number>
  errorSequences: ErrorSequence[]
  anomalyDetections: Anomaly[]
  updatedAt: string
}

interface StatisticsProps {
  stats: Statistics | null
  loading: boolean
  onRefresh: () => void
}

export default function Statistics({ stats, loading, onRefresh }: StatisticsProps) {
  if (loading) {
    return (
      <div className="flex items-center justify-center py-20">
        <RefreshCw className="h-8 w-8 animate-spin text-gray-400" />
        <span className="ml-3 text-gray-500">Loading statistics...</span>
      </div>
    )
  }

  if (!stats) {
    return (
      <div className="flex items-center justify-center py-20">
        <div className="text-center">
          <p className="text-gray-500 mb-4">Failed to load statistics.</p>
          <button
            onClick={onRefresh}
            className="inline-flex items-center px-4 py-2 border border-gray-300 shadow-sm text-sm font-medium rounded-md text-gray-700 bg-white hover:bg-gray-50"
          >
            <RefreshCw className="h-4 w-4 mr-2" />
            Retry
          </button>
        </div>
      </div>
    )
  }

  const logTypeData = Object.entries(stats.logTypeCounts).map(([name, value]) => ({ name, count: value }))
  const serviceDurationData = Object.entries(stats.serviceDurations).map(([name, value]) => ({ name, duration: value }))
  const serviceCallData = Object.entries(stats.serviceCallCounts).map(([name, value]) => ({ name, count: value }))

  return (
    <div className="space-y-8">
      <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-8">
        <Card>
          <CardHeader>
            <CardTitle>Log Types</CardTitle>
          </CardHeader>
          <CardContent>
            <ResponsiveContainer width="100%" height={300}>
              <BarChart data={logTypeData}>
                <CartesianGrid strokeDasharray="3 3" />
                <XAxis dataKey="name" />
                <YAxis />
                <Tooltip />
                <Legend />
                <Bar dataKey="count" fill="#8884d8" />
              </BarChart>
            </ResponsiveContainer>
          </CardContent>
        </Card>
        
        <Card>
          <CardHeader>
            <CardTitle>Service Durations (avg ms)</CardTitle>
          </CardHeader>
          <CardContent>
            <ResponsiveContainer width="100%" height={300}>
              <BarChart data={serviceDurationData}>
                <CartesianGrid strokeDasharray="3 3" />
                <XAxis dataKey="name" />
                <YAxis />
                <Tooltip />
                <Legend />
                <Bar dataKey="duration" fill="#82ca9d" />
              </BarChart>
            </ResponsiveContainer>
          </CardContent>
        </Card>
        
        <Card>
          <CardHeader>
            <CardTitle>Service Calls</CardTitle>
          </CardHeader>
          <CardContent>
            <ResponsiveContainer width="100%" height={300}>
              <BarChart data={serviceCallData}>
                <CartesianGrid strokeDasharray="3 3" />
                <XAxis dataKey="name" />
                <YAxis />
                <Tooltip />
                <Legend />
                <Bar dataKey="count" fill="#ffc658" />
              </BarChart>
            </ResponsiveContainer>
          </CardContent>
        </Card>
      </div>
      
      <div className="grid grid-cols-1 lg:grid-cols-2 gap-8">
        <Card>
          <CardHeader>
            <CardTitle>Error Sequences</CardTitle>
          </CardHeader>
          <CardContent>
            <div className="space-y-2">
              {stats.errorSequences && stats.errorSequences.length > 0 ? (
                stats.errorSequences.map((seq, index) => (
                  <div key={index} className="p-2 bg-gray-50 rounded">
                    <p className="font-mono text-sm">
                      {seq.service} - {new Date(seq.startTime).toLocaleString()}
                    </p>
                    <p className="text-xs text-gray-500">Count: {seq.count}</p>
                  </div>
                ))
              ) : (
                <p className="text-gray-500">No error sequences detected</p>
              )}
            </div>
          </CardContent>
        </Card>

        <Card>
          <CardHeader>
            <CardTitle>Anomaly Detections</CardTitle>
          </CardHeader>
          <CardContent>
            <div className="space-y-2">
              {stats.anomalyDetections && stats.anomalyDetections.length > 0 ? (
                stats.anomalyDetections.map((anomaly, index) => (
                  <div key={index} className="p-2 bg-red-50 border border-red-200 rounded">
                    <p className="font-bold text-sm text-red-800">{anomaly.metricName}</p>
                    <p className="text-xs text-red-600">
                      Service: {anomaly.service} - Value: {anomaly.value} (Threshold: {anomaly.threshold})
                    </p>
                    <p className="text-xs text-gray-500">
                      {new Date(anomaly.timestamp).toLocaleString()}
                    </p>
                  </div>
                ))
              ) : (
                <p className="text-gray-500">No anomalies detected</p>
              )}
            </div>
          </CardContent>
        </Card>
      </div>
    </div>
  )
}