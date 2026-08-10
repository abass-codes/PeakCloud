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
  folder_id?: string;
  name: string;
  content_type: string;
  size_bytes: number;
  etag?: string;
  created_at: string;
  updated_at: string;
};

export type Folder = {
  id: string;
  parent_id?: string;
  name: string;
  created_at: string;
  updated_at: string;
};

export type DriveContents = {
  folder_id?: string;
  breadcrumbs: Folder[];
  folders: Folder[];
  files: StoredFile[];
};

export type BulkItem = {
  type: "file" | "folder";
  id: string;
};

async function request<T>(
  path: string,
  options: RequestInit = {},
): Promise<T> {
  const response = await fetch(`${API_URL}${path}`, {
    ...options,
    credentials: "include",
    headers: {
      ...options.headers,
    },
  });

  if (!response.ok) {
    let message = "Request failed";

    try {
      const body = (await response.json()) as { error?: string };
      message = body.error ?? message;
    } catch {}

    throw new Error(message);
  }

  if (response.status === 204) {
    return undefined as T;
  }

  return (await response.json()) as T;
}

export async function register(
  email: string,
  displayName: string,
  password: string,
) {
  return request<{ user: User }>("/api/v1/auth/register", {
    method: "POST",
    headers: {
      "Content-Type": "application/json",
    },
    body: JSON.stringify({
      email,
      display_name: displayName,
      password,
    }),
  });
}

export async function login(email: string, password: string) {
  return request<{ user: User }>("/api/v1/auth/login", {
    method: "POST",
    headers: {
      "Content-Type": "application/json",
    },
    body: JSON.stringify({ email, password }),
  });
}

export async function logout() {
  return request<void>("/api/v1/auth/logout", {
    method: "POST",
  });
}

export async function getCurrentUser() {
  const result = await request<{ user: User }>("/api/v1/me");
  return result.user;
}

export async function getDrive(folderId?: string) {
  const query = folderId
    ? `?folder_id=${encodeURIComponent(folderId)}`
    : "";

  return request<DriveContents>(`/api/v1/drive${query}`);
}

export async function createFolder(
  name: string,
  parentId?: string,
) {
  const result = await request<{ folder: Folder }>("/api/v1/folders", {
    method: "POST",
    headers: {
      "Content-Type": "application/json",
    },
    body: JSON.stringify({
      name,
      parent_id: parentId ?? null,
    }),
  });

  return result.folder;
}

export async function renameFolder(id: string, name: string) {
  const result = await request<{ folder: Folder }>(
    `/api/v1/folders/${id}/name`,
    {
      method: "PATCH",
      headers: {
        "Content-Type": "application/json",
      },
      body: JSON.stringify({ name }),
    },
  );

  return result.folder;
}

export async function moveFolder(
  id: string,
  parentId?: string,
) {
  const result = await request<{ folder: Folder }>(
    `/api/v1/folders/${id}/location`,
    {
      method: "PATCH",
      headers: {
        "Content-Type": "application/json",
      },
      body: JSON.stringify({
        parent_id: parentId ?? null,
      }),
    },
  );

  return result.folder;
}

export async function deleteFolder(id: string) {
  return request<void>(`/api/v1/folders/${id}`, {
    method: "DELETE",
  });
}

export async function uploadFile(
  file: File,
  folderId?: string,
) {
  const form = new FormData();
  form.append("file", file);

  if (folderId) {
    form.append("folder_id", folderId);
  }

  const result = await request<{ file: StoredFile }>("/api/v1/files", {
    method: "POST",
    body: form,
  });

  return result.file;
}

export async function renameFile(id: string, name: string) {
  const result = await request<{ file: StoredFile }>(
    `/api/v1/files/${id}/name`,
    {
      method: "PATCH",
      headers: {
        "Content-Type": "application/json",
      },
      body: JSON.stringify({ name }),
    },
  );

  return result.file;
}

export async function moveFile(
  id: string,
  folderId?: string,
) {
  const result = await request<{ file: StoredFile }>(
    `/api/v1/files/${id}/location`,
    {
      method: "PATCH",
      headers: {
        "Content-Type": "application/json",
      },
      body: JSON.stringify({
        folder_id: folderId ?? null,
      }),
    },
  );

  return result.file;
}

export async function copyFile(
  id: string,
  folderId?: string,
) {
  const result = await request<{ file: StoredFile }>(
    `/api/v1/files/${id}/copy`,
    {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
      },
      body: JSON.stringify({
        folder_id: folderId ?? null,
      }),
    },
  );

  return result.file;
}

export async function deleteFile(id: string) {
  return request<void>(`/api/v1/files/${id}`, {
    method: "DELETE",
  });
}

export async function downloadFile(
  id: string,
  filename: string,
) {
  const response = await fetch(
    `${API_URL}/api/v1/files/${id}/download`,
    {
      credentials: "include",
    },
  );

  if (!response.ok) {
    throw new Error("Unable to download file");
  }

  const blob = await response.blob();
  const url = URL.createObjectURL(blob);
  const anchor = document.createElement("a");

  anchor.href = url;
  anchor.download = filename;
  document.body.appendChild(anchor);
  anchor.click();
  anchor.remove();

  URL.revokeObjectURL(url);
}

export async function bulkMove(
  items: BulkItem[],
  folderId?: string,
) {
  return request<void>("/api/v1/drive/bulk/move", {
    method: "POST",
    headers: {
      "Content-Type": "application/json",
    },
    body: JSON.stringify({
      items,
      folder_id: folderId ?? null,
    }),
  });
}

export async function bulkDelete(items: BulkItem[]) {
  return request<void>("/api/v1/drive/bulk/delete", {
    method: "POST",
    headers: {
      "Content-Type": "application/json",
    },
    body: JSON.stringify({ items }),
  });
}

export async function bulkDownload(fileIds: string[]) {
  const response = await fetch(
    `${API_URL}/api/v1/drive/bulk/download`,
    {
      method: "POST",
      credentials: "include",
      headers: {
        "Content-Type": "application/json",
      },
      body: JSON.stringify({
        file_ids: fileIds,
      }),
    },
  );

  if (!response.ok) {
    throw new Error("Unable to download selected files");
  }

  const blob = await response.blob();
  const url = URL.createObjectURL(blob);
  const anchor = document.createElement("a");

  anchor.href = url;
  anchor.download = "peakcloud-files.zip";
  document.body.appendChild(anchor);
  anchor.click();
  anchor.remove();

  URL.revokeObjectURL(url);
}
