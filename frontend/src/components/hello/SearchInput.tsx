interface SearchInputProps {
  value: string;
  onChange: (e: React.ChangeEvent<HTMLInputElement>) => void;
  placeholder?: string;
}

export function SearchInput({ value, onChange, placeholder }: SearchInputProps) {
  return (
    <input
      type="text"
      className="input input-bordered w-full"
      value={value}
      onChange={onChange}
      placeholder={placeholder}
    />
  );
}
