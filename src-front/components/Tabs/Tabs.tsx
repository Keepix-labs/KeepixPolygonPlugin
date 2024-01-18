import React, { useState, Children, cloneElement } from 'react';
import './Tabs.scss';

// Tab Component
export function Tab({ children }) {
  // This component is a placeholder for the tab content.
  return <>{children}</>;
}

// Tabs Component
export function Tabs({ children }) {
  const [activeTab, setActiveTab] = useState(0);
  const tabs = Children.toArray(children);

  return (
    <div className="tabs-container">
      <ul className="tab-list">
        {tabs.map((tab, index) => (
          <li
            key={index}
            className={`tab-list-item ${index === activeTab ? 'active' : ''}`}
            onClick={() => setActiveTab(index)}
          >
            {tab.props.label}
          </li>
        ))}
      </ul>

      <div className="tab-content">
        {cloneElement(tabs[activeTab])}
      </div>
    </div>
  );
}
