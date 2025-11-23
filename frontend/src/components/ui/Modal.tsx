import React, { ReactNode, useEffect } from "react";
import { createPortal } from "react-dom";

export type ModalProps = {
  isOpen: boolean;
  title?: string;
  children: ReactNode;
  footer?: ReactNode;
  onClose: () => void;
  closeOnOverlayClick?: boolean;
};

const ModalContent: React.FC<Omit<ModalProps, "isOpen">> = ({ title, children, footer, onClose, closeOnOverlayClick = true }) => {
  const handleOverlayClick = () => {
    if (closeOnOverlayClick) onClose();
  };

  return (
    <div className="fixed inset-0 z-50 flex items-center justify-center" role="dialog" aria-modal="true" aria-label={title}>
      <div className="absolute inset-0 bg-black/50" data-testid="modal-overlay" onClick={handleOverlayClick} />
      <div className="relative z-10 w-full max-w-lg rounded-lg bg-white shadow-xl">
        <div className="flex items-start justify-between border-b border-gray-200 px-4 py-3">
          <div>
            {title && <h2 className="text-lg font-semibold text-gray-900" data-testid="modal-title">{title}</h2>}
          </div>
          <button
            aria-label="閉じる"
            className="text-gray-500 hover:text-gray-700"
            onClick={onClose}
            type="button"
          >
            ×
          </button>
        </div>
        <div className="px-4 py-3 text-gray-800">{children}</div>
        {footer && <div className="border-t border-gray-200 bg-gray-50 px-4 py-3">{footer}</div>}
      </div>
    </div>
  );
};

export const Modal: React.FC<ModalProps> = (props) => {
  const { isOpen } = props;
  const [mounted, setMounted] = React.useState(false);

  useEffect(() => {
    setMounted(true);
    return () => setMounted(false);
  }, []);

  if (!isOpen || !mounted) return null;

  const portalTarget = typeof document !== "undefined" ? document.body : null;
  if (!portalTarget) return null;

  return createPortal(<ModalContent {...props} />, portalTarget);
};

export default Modal;
