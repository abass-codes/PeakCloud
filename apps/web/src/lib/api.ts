export const API_URL =
  process.env.NEXT_PUBLIC_API_URL ?? "http://localhost:8080";

export type User = {
  id: string;
  email: string;
  display_name: string;
  created_at: string;
  updated_at: string;
};

type UserResponse = {
  user: User;
};

type ErrorResponse = {
  error?: string;
};

async function parseResponse<T>(response: Response): Promise<T> {
  if (!response.ok) {
    let message = "Something went wrong";

    try {
      const body = (await response.json()) as ErrorResponse;
      message = body.error ?? message;
    } catch {
      // Keep the generic message when the response has no JSON body.
    }

    throw new Error(message);
  }

  return response.json() as Promise<T>;
}

export async function register(input: {
  email: string;
  display_name: string;
  password: string;
}): Promise<UserResponse> {
  const response = await fetch(`${API_URL}/api/v1/auth/register`, {
    method: "POST",
    credentials: "include",
    headers: {
      "Content-Type": "application/json",
    },
    body: JSON.stringify(input),
  });

  return parseResponse<UserResponse>(response);
}

export async function login(input: {
  email: string;
  password: string;
}): Promise<UserResponse> {
  const response = await fetch(`${API_URL}/api/v1/auth/login`, {
    method: "POST",
    credentials: "include",
    headers: {
      "Content-Type": "application/json",
    },
    body: JSON.stringify(input),
  });

  return parseResponse<UserResponse>(response);
}

export async function getMe(): Promise<UserResponse> {
  const response = await fetch(`${API_URL}/api/v1/me`, {
    credentials: "include",
  });

  return parseResponse<UserResponse>(response);
}

export async function logout(): Promise<void> {
  const response = await fetch(`${API_URL}/api/v1/auth/logout`, {
    method: "POST",
    credentials: "include",
  });

  if (!response.ok) {
    throw new Error("Unable to logout");
  }
}
