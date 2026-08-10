const API_URL =
  process.env.NEXT_PUBLIC_API_URL ?? "http://localhost:8080";

export type User = {
  id: string;
  email: string;
  display_name: string;
  created_at: string;
  updated_at: string;
};

export type StoredFile = {
  id: string;
  name: string;
  content_type: string;
  size_bytes: number;
  etag?: string;
  created_at: string;
  updated_at: string;
};

type UserResponse = {
  user: User;
};

type FileResponse = {
  file: StoredFile;
};

type FilesResponse = {
  files: StoredFile[];
};

async function parseError(response: Response): Promise<string> {
  try {
    const body = (await response.json()) as { error?: string };
    return body.error ?? "Request failed";
  } catch {
    return "Request failed";
  }
}

export async function register(
  email: string,
  displayName: string,
  password: string,
): Promise<User> {
  const response = await fetch(`${API_URL}/api/v1/auth/register`, {
    method: "POST",
    credentials: "include",
    headers: {
      "Content-Type": "application/json",
    },
    body: JSON.stringify({
      email,
      display_name: displayName,
      password,
    }),
  });

  if (!response.ok) {
    throw new Error(await parseError(response));
  }

  const body = (await response.json()) as UserResponse;
  return body.user;
}

export async function login(
  email: string,
  password: string,
): Promise<User> {
  const response = await fetch(`${API_URL}/api/v1/auth/login`, {
    method: "POST",
    credentials: "include",
    headers: {
      "Content-Type": "application/json",
    },
    body: JSON.stringify({
      email,
      password,
    }),
  });

  if (!response.ok) {
    throw new Error(await parseError(response));
  }

  const body = (await response.json()) as UserResponse;
  return body.user;
}

export async function logout(): Promise<void> {
  const response = await fetch(`${API_URL}/api/v1/auth/logout`, {
    method: "POST",
    credentials: "include",
  });

  if (!response.ok && response.status !== 204) {
    throw new Error(await parseError(response));
  }
}

export async function getCurrentUser(): Promise<User> {
  const response = await fetch(`${API_URL}/api/v1/me`, {
    credentials: "include",
  });

  if (!response.ok) {
    throw new Error(await parseError(response));
  }

  const body = (await response.json()) as UserResponse;
  return body.user;
}

export async function listFiles(): Promise<StoredFile[]> {
  const response = await fetch(`${API_URL}/api/v1/files`, {
    credentials: "include",
  });

  if (!response.ok) {
    throw new Error(await parseError(response));
  }

  const body = (await response.json()) as FilesResponse;
  return body.files;
}

export async function uploadFile(file: File): Promise<StoredFile> {
  const form = new FormData();
  form.append("file", file);

  const response = await fetch(`${API_URL}/api/v1/files`, {
    method: "POST",
    credentials: "include",
    body: form,
  });

  if (!response.ok) {
    throw new Error(await parseError(response));
  }

  const body = (await response.json()) as FileResponse;
  return body.file;
}

export async function deleteFile(id: string): Promise<void> {
  const response = await fetch(`${API_URL}/api/v1/files/${id}`, {
    method: "DELETE",
    credentials: "include",
  });

  if (!response.ok && response.status !== 204) {
    throw new Error(await parseError(response));
  }
}

export async function downloadFile(file: StoredFile): Promise<void> {
  const response = await fetch(
    `${API_URL}/api/v1/files/${file.id}/download`,
    {
      credentials: "include",
    },
  );

  if (!response.ok) {
    throw new Error(await parseError(response));
  }

  const blob = await response.blob();
  const url = window.URL.createObjectURL(blob);

  const anchor = document.createElement("a");
  anchor.href = url;
  anchor.download = file.name;
  document.body.appendChild(anchor);
  anchor.click();
  anchor.remove();

  window.URL.revokeObjectURL(url);
}
