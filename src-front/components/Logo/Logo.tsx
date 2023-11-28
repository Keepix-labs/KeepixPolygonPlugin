import "./Logo.scss";

const Logo = ({ text = true, width = undefined }: any) => {
  const y = 90.5;
  const x = 156.5;

  const cubes = [
    { y: 0, x: 4 },
    { y: 1, x: 3 },
    { y: 2, x: 2 },
    { y: 3, x: 1 },
    { y: 4, x: 0 },
    { y: 3, x: 3, secondary: true },
    { y: 3, x: 5, secondary: true },
    { y: 4, x: 4, secondary: true },
    { y: 5, x: 3, secondary: true },
    { y: 7, x: 3 },
    { y: 3, x: 7 },
  ];

  return (
    <div className="Logo-logo">
      <svg className="Logo-logoSigle" style={{ width: width }} viewBox="0 0 1437 1020">
        {cubes.map((cube, index) => (
          <use
            key={index}
            xlinkHref={cube.secondary ? `#cube-b-secondary` : `#cube-b`}
            className={`cube-${index}`}
            y={cube.y * y}
            x={cube.x * x}
            style={{ animationDelay: `.${index}s` }}
          />
        ))}
      </svg>
      {text && (
        <svg className="Logo-logoTxt" viewBox="0 0 431.8 114.9">
          <use xlinkHref="#logo-txt" />
        </svg>
      )}
    </div>
  );
};

export default Logo;
