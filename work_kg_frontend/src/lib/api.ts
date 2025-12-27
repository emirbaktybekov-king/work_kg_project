const API_URL = process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8080/api';

export interface User {
  id: number;
  telegram_id: number;
  username: string;
  first_name: string;
  last_name: string;
  phone: string;
  city: string;
  specialty: string;
  experience: string;
  role: string;
  created_at: string;
}

export interface AdminUser {
  id: number;
  email: string;
  name: string;
  role: string;
  created_at: string;
}

export interface Job {
  id: number;
  title: string;
  description: string;
  category: string;
  subcategory: string;
  city: string;
  salary: string;
  phone: string;
  company: string;
  is_active: boolean;
  created_by: number;
  source: string;
  created_at: string;
}

export interface Resume {
  id: number;
  telegram_id: number;
  username: string;
  name: string;
  phone: string;
  city: string;
  specialty: string;
  experience: string;
  created_at: string;
  updated_at: string;
}

export interface Stats {
  total_jobs: number;
  active_jobs: number;
  total_users: number;
  total_resumes: number;
  today_jobs: number;
  today_users: number;
  today_resumes: number;
}

export interface LoginResponse {
  token: string;
  user: AdminUser;
}

class ApiClient {
  private token: string | null = null;

  setToken(token: string | null) {
    this.token = token;
    if (typeof window !== 'undefined') {
      if (token) {
        localStorage.setItem('token', token);
      } else {
        localStorage.removeItem('token');
      }
    }
  }

  getToken(): string | null {
    if (this.token) return this.token;
    if (typeof window !== 'undefined') {
      this.token = localStorage.getItem('token');
    }
    return this.token;
  }

  private async request<T>(endpoint: string, options: RequestInit = {}): Promise<T> {
    const token = this.getToken();
    const headers: Record<string, string> = {
      'Content-Type': 'application/json',
    };

    if (token) {
      headers['Authorization'] = `Bearer ${token}`;
    }

    const response = await fetch(`${API_URL}${endpoint}`, {
      ...options,
      headers,
    });

    if (!response.ok) {
      const error = await response.text();
      throw new Error(error || 'Request failed');
    }

    if (response.status === 204) {
      return {} as T;
    }

    return response.json();
  }

  // Auth
  async login(email: string, password: string): Promise<LoginResponse> {
    const response = await this.request<LoginResponse>('/auth/login', {
      method: 'POST',
      body: JSON.stringify({ email, password }),
    });
    this.setToken(response.token);
    return response;
  }

  async getMe(): Promise<AdminUser> {
    return this.request<AdminUser>('/auth/me');
  }

  logout() {
    this.setToken(null);
    if (typeof window !== 'undefined') {
      localStorage.removeItem('user');
      localStorage.removeItem('role');
    }
  }

  // Jobs
  async getJobs(): Promise<Job[]> {
    return this.request<Job[]>('/jobs');
  }

  async createJob(job: Partial<Job>): Promise<Job> {
    return this.request<Job>('/jobs', {
      method: 'POST',
      body: JSON.stringify(job),
    });
  }

  async updateJob(id: number, job: Partial<Job>): Promise<Job> {
    return this.request<Job>(`/jobs/${id}`, {
      method: 'PUT',
      body: JSON.stringify(job),
    });
  }

  async deleteJob(id: number): Promise<void> {
    return this.request<void>(`/jobs/${id}`, {
      method: 'DELETE',
    });
  }

  // Users
  async getUsers(): Promise<User[]> {
    return this.request<User[]>('/users');
  }

  // Resumes
  async getResumes(): Promise<Resume[]> {
    return this.request<Resume[]>('/resumes');
  }

  // Stats
  async getStats(): Promise<Stats> {
    return this.request<Stats>('/stats');
  }
}

export const api = new ApiClient();
