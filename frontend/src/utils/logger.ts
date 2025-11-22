/**
 * Frontend Logger Utility
 * Provides structured logging for development and production environments
 */

export enum LogLevel {
  DEBUG = 'DEBUG',
  INFO = 'INFO',
  WARN = 'WARN',
  ERROR = 'ERROR',
}

interface LogEntry {
  timestamp: string
  level: LogLevel
  message: string
  context?: string
  data?: unknown
  error?: Error
}

class Logger {
  private isDevelopment: boolean
  private enableConsole: boolean
  private minLevel: LogLevel

  constructor() {
    this.isDevelopment = import.meta.env.DEV
    this.enableConsole = true
    this.minLevel = this.isDevelopment ? LogLevel.DEBUG : LogLevel.INFO
  }

  private shouldLog(level: LogLevel): boolean {
    const levels = [LogLevel.DEBUG, LogLevel.INFO, LogLevel.WARN, LogLevel.ERROR]
    return levels.indexOf(level) >= levels.indexOf(this.minLevel)
  }

  private formatMessage(entry: LogEntry): string {
    const parts = [`[${entry.timestamp}]`, `[${entry.level}]`]

    if (entry.context) {
      parts.push(`[${entry.context}]`)
    }

    parts.push(entry.message)

    return parts.join(' ')
  }

  private createLogEntry(
    level: LogLevel,
    message: string,
    context?: string,
    data?: unknown,
    error?: Error,
  ): LogEntry {
    return {
      timestamp: new Date().toISOString(),
      level,
      message,
      context,
      data,
      error,
    }
  }

  private log(entry: LogEntry): void {
    if (!this.shouldLog(entry.level)) {
      return
    }

    const formattedMessage = this.formatMessage(entry)

    if (this.enableConsole) {
      switch (entry.level) {
        case LogLevel.DEBUG:
          console.debug(formattedMessage, entry.data || '')
          break
        case LogLevel.INFO:
          console.info(formattedMessage, entry.data || '')
          break
        case LogLevel.WARN:
          console.warn(formattedMessage, entry.data || '')
          break
        case LogLevel.ERROR:
          console.error(formattedMessage, entry.data || '', entry.error || '')
          break
      }
    }

    // In production, you could send logs to a remote service here
    if (!this.isDevelopment && entry.level === LogLevel.ERROR) {
      this.sendToRemote(entry)
    }
  }

  private sendToRemote(entry: LogEntry): void {
    // Placeholder for remote logging service (e.g., Sentry, LogRocket)
    // Example: Sentry.captureException(entry.error)
    if (this.isDevelopment) {
      console.log('[Remote logging would send:]', entry)
    }
  }

  debug(message: string, context?: string, data?: unknown): void {
    this.log(this.createLogEntry(LogLevel.DEBUG, message, context, data))
  }

  info(message: string, context?: string, data?: unknown): void {
    this.log(this.createLogEntry(LogLevel.INFO, message, context, data))
  }

  warn(message: string, context?: string, data?: unknown): void {
    this.log(this.createLogEntry(LogLevel.WARN, message, context, data))
  }

  error(message: string, context?: string, error?: Error, data?: unknown): void {
    this.log(this.createLogEntry(LogLevel.ERROR, message, context, data, error))
  }

  // HTTP request logging helpers
  logRequest(method: string, url: string, data?: unknown): void {
    this.debug(`${method} ${url}`, 'HTTP', data)
  }

  logResponse(method: string, url: string, status: number, data?: unknown): void {
    if (status >= 400) {
      this.warn(`${method} ${url} - ${status}`, 'HTTP', data)
    } else {
      this.debug(`${method} ${url} - ${status}`, 'HTTP', data)
    }
  }

  logError(method: string, url: string, error: Error): void {
    this.error(`${method} ${url} failed`, 'HTTP', error)
  }

  // Store action logging
  logStoreAction(store: string, action: string, payload?: unknown): void {
    this.debug(`${store}.${action}`, 'STORE', payload)
  }

  logStoreError(store: string, action: string, error: Error): void {
    this.error(`${store}.${action} failed`, 'STORE', error)
  }

  // Component lifecycle logging
  logComponentMount(component: string): void {
    this.debug(`Component mounted`, component)
  }

  logComponentUnmount(component: string): void {
    this.debug(`Component unmounted`, component)
  }
}

// Export singleton instance
export const logger = new Logger()

// Export for testing
export default logger
