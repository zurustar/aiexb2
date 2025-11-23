import React, { ChangeEvent, forwardRef } from "react";

export type DatePickerProps = {
  label?: string;
  value?: string;
  onChange?: (value: string) => void;
  min?: string;
  max?: string;
  required?: boolean;
  disabled?: boolean;
  helperText?: string;
  error?: string;
  name?: string;
};

export const DatePicker = forwardRef<HTMLInputElement, DatePickerProps>(
  ({ label, value, onChange, min, max, required, disabled, helperText, error, name }, ref) => {
    const handleChange = (event: ChangeEvent<HTMLInputElement>) => {
      onChange?.(event.target.value);
    };

    return (
      <div className="flex flex-col gap-1">
        {label && (
          <label className="text-sm font-medium text-gray-800" htmlFor={name}>
            {label}
            {required && <span className="text-red-500">*</span>}
          </label>
        )}
        <input
          ref={ref}
          id={name}
          name={name}
          data-testid={name}
          type="datetime-local"
          value={value}
          onChange={handleChange}
          min={min}
          max={max}
          required={required}
          disabled={disabled}
          className={`rounded-md border px-3 py-2 shadow-sm focus:border-blue-500 focus:outline-none focus:ring-2 focus:ring-blue-200 ${error ? "border-red-500" : "border-gray-300"
            }`}
        />
        {helperText && !error && <p className="text-xs text-gray-500">{helperText}</p>}
        {error && <p className="text-xs text-red-600" role="alert">{error}</p>}
      </div>
    );
  }
);

DatePicker.displayName = "DatePicker";

export default DatePicker;
