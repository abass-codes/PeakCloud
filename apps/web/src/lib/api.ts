const API_URL = "/api/backend";

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

async function request<T>(path: string, options: RequestInit = {}): Promise<T> {
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
  const query = folderId ? `?folder_id=${encodeURIComponent(folderId)}` : "";

  return request<DriveContents>(`/api/v1/drive${query}`);
}

export async function createFolder(name: string, parentId?: string) {
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

export async function moveFolder(id: string, parentId?: string) {
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

export async function uploadFile(file: File, folderId?: string) {
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

export async function moveFile(id: string, folderId?: string) {
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

export async function copyFile(id: string, folderId?: string) {
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

export async function downloadFile(id: string, filename: string) {
  const response = await fetch(`${API_URL}/api/v1/files/${id}/download`, {
    credentials: "include",
  });

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

export async function bulkMove(items: BulkItem[], folderId?: string) {
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
  const response = await fetch(`${API_URL}/api/v1/drive/bulk/download`, {
    method: "POST",
    credentials: "include",
    headers: {
      "Content-Type": "application/json",
    },
    body: JSON.stringify({
      file_ids: fileIds,
    }),
  });

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

export type PreviewKind =
  "image" | "pdf" | "text" | "code" | "video" | "audio" | "unsupported";

export type FilePreviewInfo = {
  kind: PreviewKind;
  previewable: boolean;
  inline: boolean;
};

export type FilePreviewResponse = {
  file: StoredFile;
  preview: FilePreviewInfo;
};

export async function getFilePreview(
  fileId: string,
): Promise<FilePreviewResponse> {
  const response = await fetch(`${API_URL}/api/v1/files/${fileId}/preview`, {
    credentials: "include",
    cache: "no-store",
  });

  if (!response.ok) {
    throw new Error("Unable to load file preview");
  }

  return response.json();
}

export async function getFilePreviewBlob(fileId: string): Promise<Blob> {
  const response = await fetch(`${API_URL}/api/v1/files/${fileId}/content`, {
    credentials: "include",
    cache: "no-store",
  });

  if (!response.ok) {
    throw new Error("Unable to load file content");
  }

  return response.blob();
}

export type ShareResourceType = "file" | "folder";
export type SharePermission = "viewer" | "editor";

export type ResourceShare = {
  id: string;
  owner_id: string;
  recipient_id: string;
  recipient_email?: string;
  resource_type: ShareResourceType;
  resource_id: string;
  resource_name: string;
  permission: SharePermission;
  allow_download: boolean;
  created_at: string;
  updated_at: string;
};

export type PublicShareLink = {
  id: string;
  resource_type: ShareResourceType;
  resource_id: string;
  resource_name: string;
  permission: SharePermission;
  allow_download: boolean;
  password_set: boolean;
  expires_at?: string;
  revoked_at?: string;
  created_at: string;
  updated_at: string;
};

export async function createResourceShare(input: {
  recipient_email: string;
  resource_type: ShareResourceType;
  resource_id: string;
  permission: SharePermission;
  allow_download: boolean;
}) {
  const result = await request<{ share: ResourceShare }>("/api/v1/shares", {
    method: "POST",
    headers: {
      "Content-Type": "application/json",
    },
    body: JSON.stringify(input),
  });

  return result.share;
}

export async function getResourceShares() {
  const result = await request<{ shares: ResourceShare[] }>("/api/v1/shares");
  return result.shares;
}

export async function getSharedWithMe() {
  const result = await request<{ shares: ResourceShare[] }>(
    "/api/v1/shared-with-me",
  );

  return result.shares;
}

export async function updateResourceShare(
  id: string,
  input: {
    permission: SharePermission;
    allow_download: boolean;
  },
) {
  return request<void>(`/api/v1/shares/${id}`, {
    method: "PATCH",
    headers: {
      "Content-Type": "application/json",
    },
    body: JSON.stringify(input),
  });
}

export async function deleteResourceShare(id: string) {
  return request<void>(`/api/v1/shares/${id}`, {
    method: "DELETE",
  });
}

export async function createPublicShareLink(input: {
  resource_type: ShareResourceType;
  resource_id: string;
  permission: SharePermission;
  allow_download: boolean;
  password?: string;
  expires_at?: string;
}) {
  const result = await request<{
    link: PublicShareLink & { token: string };
  }>("/api/v1/public-links", {
    method: "POST",
    headers: {
      "Content-Type": "application/json",
    },
    body: JSON.stringify(input),
  });

  return result.link;
}

export async function getPublicShareLinks() {
  const result = await request<{ links: PublicShareLink[] }>(
    "/api/v1/public-links",
  );

  return result.links;
}

export async function revokePublicShareLink(id: string) {
  return request<void>(`/api/v1/public-links/${id}`, {
    method: "DELETE",
  });
}

export async function resolvePublicShare(
  token: string,
  password = "",
): Promise<PublicShareLink> {
  const response = await fetch(`${API_URL}/api/v1/public/shares/${token}`, {
    method: "POST",
    headers: {
      "Content-Type": "application/json",
    },
    body: JSON.stringify({ password }),
    cache: "no-store",
  });

  const body = (await response.json()) as {
    link?: PublicShareLink;
    error?: string;
  };

  if (!response.ok || !body.link) {
    throw new Error(body.error ?? "Unable to open shared resource");
  }

  return body.link;
}

export function publicShareContentUrl(token: string, password = "") {
  const query = password ? `?password=${encodeURIComponent(password)}` : "";

  return `${API_URL}/api/v1/public/shares/${token}/content${query}`;
}

export function publicShareDownloadUrl(token: string, password = "") {
  const query = password ? `?password=${encodeURIComponent(password)}` : "";

  return `${API_URL}/api/v1/public/shares/${token}/download${query}`;
}

export type FileVersion = {
  id: string;
  file_id: string;
  version_number: number;
  size_bytes: number;
  content_type: string;
  etag?: string;
  created_by: string;
  created_at: string;
};

export async function getFileVersions(fileId: string) {
  const result = await request<{ versions: FileVersion[] }>(
    `/api/v1/files/${fileId}/versions`,
  );

  return result.versions;
}

export async function getFileVersion(
  fileId: string,
  versionNumber: number,
) {
  const result = await request<{ version: FileVersion }>(
    `/api/v1/files/${fileId}/versions/${versionNumber}`,
  );

  return result.version;
}

export async function uploadFileVersion(
  fileId: string,
  file: File,
) {
  const form = new FormData();
  form.append("file", file);

  const result = await request<{ version: FileVersion }>(
    `/api/v1/files/${fileId}/versions`,
    {
      method: "POST",
      body: form,
    },
  );

  return result.version;
}

export async function getFileVersionContent(
  fileId: string,
  versionNumber: number,
): Promise<Blob> {
  const response = await fetch(
    `${API_URL}/api/v1/files/${fileId}/versions/${versionNumber}/content`,
    {
      credentials: "include",
      cache: "no-store",
    },
  );

  if (!response.ok) {
    throw new Error("Unable to load file version content");
  }

  return response.blob();
}

export async function downloadFileVersion(
  fileId: string,
  versionNumber: number,
  filename: string,
) {
  const response = await fetch(
    `${API_URL}/api/v1/files/${fileId}/versions/${versionNumber}/download`,
    {
      credentials: "include",
    },
  );

  if (!response.ok) {
    throw new Error("Unable to download file version");
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

export async function restoreFileVersion(
  fileId: string,
  versionNumber: number,
) {
  const result = await request<{ version: FileVersion }>(
    `/api/v1/files/${fileId}/versions/${versionNumber}/restore`,
    {
      method: "POST",
    },
  );

  return result.version;
}

export type TrashResourceType = "file" | "folder";

export type TrashItem = {
  id: string;
  resource_type: TrashResourceType;
  name: string;
  content_type?: string;
  size_bytes?: number;
  deleted_at: string;
  created_at: string;
  updated_at: string;
};

export async function getTrash() {
  return request<{ items: TrashItem[] }>(
    "/api/v1/trash",
  );
}

export async function moveToTrash(
  resourceType: TrashResourceType,
  resourceId: string,
) {
  return request<void>(
    `/api/v1/trash/${resourceType}/${resourceId}`,
    {
      method: "POST",
    },
  );
}

export async function restoreTrashItem(
  resourceType: TrashResourceType,
  resourceId: string,
) {
  return request<void>(
    `/api/v1/trash/${resourceType}/${resourceId}/restore`,
    {
      method: "POST",
    },
  );
}

export async function permanentlyDeleteTrashItem(
  resourceType: TrashResourceType,
  resourceId: string,
) {
  return request<void>(
    `/api/v1/trash/${resourceType}/${resourceId}`,
    {
      method: "DELETE",
    },
  );
}
