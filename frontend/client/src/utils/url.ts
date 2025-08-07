export class UrlUtils {
  static getRoomIdFromUrl(): string | null {
    const params = new URLSearchParams(window.location.search);
    return params.get('room');
  }

  static setRoomIdInUrl(roomId: string): void {
    const url = new URL(window.location.href);
    url.searchParams.set('room', roomId);
    window.history.pushState({}, '', url.toString());
  }

  static removeRoomIdFromUrl(): void {
    const url = new URL(window.location.href);
    url.searchParams.delete('room');
    window.history.pushState({}, '', url.toString());
  }

  static getShareUrl(roomId: string): string {
    const baseUrl = window.location.origin;
    return `${baseUrl}/editor?room=${roomId}`;
  }

  static copyToClipboard(text: string): Promise<void> {
    if (navigator.clipboard) {
      return navigator.clipboard.writeText(text);
    } else {
      // Fallback for older browsers
      const textArea = document.createElement('textarea');
      textArea.value = text;
      document.body.appendChild(textArea);
      textArea.select();
      document.execCommand('copy');
      document.body.removeChild(textArea);
      return Promise.resolve();
    }
  }

  static generateRoomId(): string {
    return `room_${Date.now()}_${Math.random().toString(36).substr(2, 9)}`;
  }
}

