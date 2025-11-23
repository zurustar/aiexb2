import React from "react";

export type ToastVariant = "info" | "success" | "warning" | "error";

export type ToastProps = {
  message: string;
  variant?: ToastVariant;
  isVisible?: boolean;
  onClose?: () => void;
  duration?: number;
};

const variantStyles: Record<ToastVariant, string> = {
  info: "bg-blue-50 text-blue-900 border-blue-200",
  success: "bg-green-50 text-green-900 border-green-200",
  warning: "bg-yellow-50 text-yellow-900 border-yellow-200",
  error: "bg-red-50 text-red-900 border-red-200",
};

export const Toast: React.FC<ToastProps> = ({
  message,
  variant = "info",
  isVisible = true,
  onClose,
  duration = 4000,
}) => {
  React.useEffect(() => {
    if (!isVisible || !onClose) return;
    const timer = setTimeout(onClose, duration);
    return () => clearTimeout(timer);
  }, [isVisible, onClose, duration]);

  if (!isVisible) return null;

  return (
    <div
      role="status"
      className={`pointer-events-auto flex items-start gap-3 rounded-md border px-4 py-3 shadow-md ${variantStyles[variant]}`}
    >
      <span className="text-lg" aria-hidden>
        {variant === "success" && "✔"}
        {variant === "error" && "⚠"}
        {variant === "warning" && "!"}
        {variant === "info" && "ℹ"}
      </span>
      <div className="flex-1 text-sm font-medium">{message}</div>
      {onClose && (
        <button type="button" onClick={onClose} aria-label="閉じる" className="text-sm font-semibold">
          ×
        </button>
      )}
    </div>
  );
};

export default Toast;
