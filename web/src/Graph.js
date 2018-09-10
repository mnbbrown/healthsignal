import React, { Component } from "react";
import {
  styler,
  Charts,
  Resizable,
  ChartContainer,
  ChartRow,
  YAxis,
  LineChart
} from "react-timeseries-charts";
import { TimeSeries } from "pondjs";
import api from "./api";
import "./Graph.css";

class Graph extends Component {
  state = {
    data: {},
    locations: []
  };

  syncData = () => {
    const { endpoint } = this.props;
    api
      .getEndpointData(endpoint)
      .then(data => {
        data = data.reduce((dataByTime, point) => {
          dataByTime[point["timestamp"]] = dataByTime[point["timestamp"]] || {};
          const location = point["location"];
          dataByTime[point["timestamp"]][location] = point;
          return dataByTime;
        }, {});

        let locations = new Set();
        let columns = new Set();
        data = Object.keys(data).map(timestamp => {
          const tick = data[timestamp];
          const locationNames = Object.keys(tick);
          return locationNames.reduce(
            (prefixed, location) => {
              locations.add(location);
              const p = Object.keys(tick[location]).reduce((p, dps) => {
                if (dps === "location" || dps === "timestamp") {
                  return p;
                }
                const header = `${location}_${dps}`;
                columns.add(header);
                return { ...p, [`${location}_${dps}`]: tick[location][dps] };
              }, {});
              return { ...prefixed, ...p };
            },
            { time: parseInt(timestamp, 10) }
          );
        });
        columns = ["time", ...columns];
        locations = [...locations];
        data = data.map(row => {
          return columns.map(c => row[c]);
        });
        data = new TimeSeries({
          name: this.props.endpoint.name,
          points: data,
          columns
        });
        this.setState({ locations, data });
      })
      .catch(e => console.error(e));
  };

  componentDidMount() {
    this.syncData();
    this.poller = setInterval(() => {
      this.syncData();
    }, 5000);
  }

  componentWillUnmount() {
    clearInterval(this.poller);
  }

  render() {
    const { data, locations } = this.state;
    if (!data || locations.length === 0) {
      return null;
    }
    const colors = ["steelblue", "green", "pink"];
    const lineStyle = new styler(
      locations.map((location, index) => {
        return {
          key: `${location}_requestTime`,
          color: colors[index],
          width: 2
        };
      })
    );
    return (
      data && (
        <div className="Graph">
          <h3>{this.props.endpoint.name}</h3>
          <Resizable>
            <ChartContainer
              timeRange={data.timerange()}
              enablePanZoom={true}
              width={800}
            >
              <ChartRow height="100">
                <YAxis
                  id="response"
                  min={0}
                  max={1200}
                  width="60"
                  label="ms"
                  type="linear"
                />
                <Charts>
                  <LineChart
                    lineStyle={lineStyle}
                    axis="response"
                    series={data}
                    columns={locations.map(l => `${l}_requestTime`)}
                  />
                </Charts>
              </ChartRow>
            </ChartContainer>
          </Resizable>
        </div>
      )
    );
  }
}

export default Graph;
