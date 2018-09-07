import React, { Component } from "react";
import {
  styler,
  Charts,
  Resizable,
  ChartContainer,
  ChartRow,
  YAxis,
  LineChart,
  EventMarker
} from "react-timeseries-charts";
import { TimeSeries } from "pondjs";
import "./Graph.css";

class Graph extends Component {
  state = {
    data: null
  };

  syncData = () => {
    const { endpoint } = this.props;
    fetch(`https://api.healthsignal.live/endpoints/${endpoint.id}/data`, {
      method: "GET"
    })
      .then(response => response.json())
      .then(data => {
        const timeseries = new TimeSeries({
          name: "1",
          columns: ["time", "requestTime", "connectionTime"],
          points: data.map(p => [
            p["timestamp"],
            p["requestTime"],
            p["connectionTime"]
          ])
        });
        this.setState({
          data: this.state.data
            ? TimeSeries.timeSeriesListMerge({
                name: "1",
                seriesList: [timeseries, this.state.data]
              })
            : timeseries
        });
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

  handleTrackerChanged = t => {
    if (t) {
      const e = this.state.data.atTime(t);
      const eventTime = new Date(
        e.begin().getTime() + (e.end().getTime() - e.begin().getTime()) / 2
      );
      const eventValue = e.get("requestTime");
      const v = `${Math.round(eventValue)}ms`;
      this.setState({ tracker: eventTime, trackerValue: v, trackerEvent: e });
    } else {
      this.setState({ tracker: null, trackerValue: null, trackerEvent: null });
    }
  };

  render() {
    const { data } = this.state;
    const lineStyle = styler([
      { key: "requestTime", color: "steelblue", width: 2 }
    ]);
    return (
      data && (
        <div className="Graph">
          <h3>{this.props.endpoint.name}</h3>
          <Resizable>
            <ChartContainer
              timeRange={data.timerange()}
              enablePanZoom={true}
              width={800}
              onTrackerChanged={this.handleTrackerChanged}
            >
              <ChartRow height="100">
                <YAxis
                  id="response"
                  min={0}
                  max={data.max("requestTime") * 2}
                  width="60"
                  label="ms"
                  type="linear"
                />
                <Charts>
                  <LineChart
                    style={lineStyle}
                    interpolation="curveLinear"
                    axis="response"
                    series={data}
                    columns={["requestTime"]}
                  />
                  <EventMarker
                    event={this.state.trackerEvent}
                    markerLabel={this.state.trackerValue}
                    markerLabelAlign="top"
                    markerLabelStyle={{ fill: "#000", stroke: "white" }}
                    type="point"
                    axis="response"
                    column="requestTime"
                    markerRadius={2}
                    markerStyle={{ fill: "black" }}
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
