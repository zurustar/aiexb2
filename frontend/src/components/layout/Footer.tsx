import React from "react";

export type FooterLink = {
  label: string;
  href: string;
};

export type FooterProps = {
  links?: FooterLink[];
  version?: string;
  companyName?: string;
};

export const Footer: React.FC<FooterProps> = ({ links = [], version, companyName = "ESMS" }) => {
  const currentYear = new Date().getFullYear();

  return (
    <footer className="flex flex-col gap-3 border-t border-gray-200 bg-white px-4 py-3 text-sm text-gray-700 md:flex-row md:items-center md:justify-between">
      <div className="flex items-center gap-2">
        <span className="font-semibold">{companyName}</span>
        <span className="text-gray-400">© {currentYear}</span>
        {version && <span className="rounded bg-gray-100 px-2 py-0.5 text-xs font-semibold">v{version}</span>}
      </div>
      <div className="flex flex-wrap items-center gap-3">
        {links.map((link) => (
          <a key={link.href} href={link.href} className="text-blue-700 hover:underline">
            {link.label}
          </a>
        ))}
      </div>
    </footer>
  );
};

export default Footer;
