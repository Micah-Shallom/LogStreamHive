"use client"

import { useState, useEffect } from "react"
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card"
import { BarChart, Bar, XAxis, YAxis, CartesianGrid, Tooltip, Legend, ResponsiveContainer } from 'recharts';

interface Statistics {
  logTypeCounts: Record<string, number>;
  serviceDurations: Record<string, number>;
  serviceCallCounts: Record<string, number>;
  errorSequences: ErrorSequence[];
  anomalyDetections: Anomaly[];
  updatedAt: string;
}

interface ErrorSequence {
  startTime: string;
  endTime: string;
  count: number;
  service: string;
}

interface Anomaly {
  timestamp: string;
  service: string;
  metricName: string;
  value: number;
  threshold: number;
}

export default function Statistics() {
  const [stats, setStats] = useState<Statistics | null>(null)
  const [loading, setLoading] = useState(true)

  const API_URL = process.env.NEXT_PUBLIC_API_URL || "http://localhost:8000"

  const fetchStatistics = async () => {
    setLoading(true)
    try {
      const response = await fetch(`${API_URL}/statistics`)
      if (!response.ok) {
        throw new Error(`HTTP error! status: ${response.status}`)
      }
      const data: Statistics = await response.json()

      if (!data.logTypeCounts || !data.serviceDurations || !data.serviceCallCounts) {
        throw new Error("Incomplete statistics data received from server.")
      }

      console.log(data)

      setStats(data as Statistics)
    } catch (error) {
      console.error("Failed to fetch statistics:", error)
      setStats(null)
    } finally {
      setLoading(false)
    }
  }

  useEffect(() => {
    fetchStatistics()
  }, [])

  if (loading) {
    return <div>Loading statistics...</div>
  }

  if (!stats) {
    return <div>Failed to load statistics.</div>
  }

  const logTypeData = Object.entries(stats.logTypeCounts).map(([name, value]) => ({ name, count: value }));
  const serviceDurationData = Object.entries(stats.serviceDurations).map(([name, value]) => ({ name, duration: value }));
  const serviceCallData = Object.entries(stats.serviceCallCounts).map(([name, value]) => ({ name, count: value }));

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
              {stats.errorSequences?.map((seq, index) => (
                <div key={index} className="p-2 bg-gray-50 rounded">
                  <p className="font-mono text-sm">
                    {seq.service} - {new Date(seq.startTime).toLocaleString()}
                  </p>
                  <p className="text-xs text-gray-500">Count: {seq.count}</p>
                </div>
              )) ?? <p>No error sequences detected</p>}
            </div>
          </CardContent>
        </Card>

        <Card>
          <CardHeader>
            <CardTitle>Anomaly Detections</CardTitle>
          </CardHeader>
          <CardContent>
            <div className="space-y-2">
              {stats.anomalyDetections?.map((anomaly, index) => (
                <div key={index} className="p-2 bg-red-50 border border-red-200 rounded">
                  <p className="font-bold text-sm text-red-800">{anomaly.metricName}</p>
                  <p className="text-xs text-red-600">
                    Service: {anomaly.service} - Value: {anomaly.value} (Threshold: {anomaly.threshold})
                  </p>
                  <p className="text-xs text-gray-500">
                    {new Date(anomaly.timestamp).toLocaleString()}
                  </p>
                </div>
              )) ?? <p>No anomalies detected</p>}
            </div>
          </CardContent>
        </Card>
      </div>
    </div>
  )
}