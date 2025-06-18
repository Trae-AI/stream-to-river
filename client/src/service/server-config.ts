// Copyright (c) 2025 Bytedance Ltd. and/or its affiliates 
// SPDX-License-Identifier: MIT 

interface ServerConfigInterface {
  domain: string;
  port?: number;
  protocol?: string;
}

class ServerConfig {
  private static instance: ServerConfig;
  private config: ServerConfigInterface;

  private constructor() {
    this.config = this.loadConfig();
  }

  public static getInstance(): ServerConfig {
    if (!ServerConfig.instance) {
      ServerConfig.instance = new ServerConfig();
    }
    return ServerConfig.instance;
  }

  private loadConfig(): ServerConfigInterface {
    const serverDomain = location.origin;

    const url = new URL(serverDomain);

    return {
      domain: url.hostname,
      port: url.port ? parseInt(url.port) : (url.protocol === 'https:' ? 443 : 80),
      protocol: url.protocol.replace(':', ''),
    };
  }

  public getDomain(): string {
    return this.config.domain;
  }

  public getPort(): number {
    return this.config.port || 80;
  }

  public getProtocol(): string {
    return this.config.protocol || 'http';
  }

  public getBaseUrl(): string {
    const port = this.getPort();
    const protocol = this.getProtocol();
    const domain = this.getDomain();

    if ((protocol === 'http' && port === 80) || (protocol === 'https' && port === 443)) {
      return `${protocol}://${domain}`;
    }

    return `${protocol}://${domain}:${port}`;
  }

  public getFullUrl(path: string = ''): string {
    const baseUrl = this.getBaseUrl();
    const cleanPath = path.startsWith('/') ? path : `/${path}`;
    return `${baseUrl}${cleanPath}`;
  }
}

export { ServerConfigInterface, ServerConfig };
export default ServerConfig;
