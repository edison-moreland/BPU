//visualization: netlistsvg-sop

// 0 0 0 0 |7'b1111110
// 0 0 0 1 |7'b0110000
// 0 0 1 0 |7'b1101101
// 0 0 1 1 |7'b1111001
// 0 1 0 0 |7'b0110011
// 0 1 0 1 |7'b1011011
// 0 1 1 0 |7'b1011111
// 0 1 1 1 |7'b1110000
// 1 0 0 0 |7'b1111111
// 1 0 0 1 |7'b1111011
// 1 0 1 0 |7'b1110111
// 1 0 1 1 |7'b0011111
// 1 1 0 0 |7'b1001110
// 1 1 0 1 |7'b0111101
// 1 1 1 0 |7'b1001111
// 1 1 1 1 |7'b1000111

module hex_decoder (
        input  logic[0:3] value,
        output logic[0:6] display
);
  always_comb begin
    unique case (value)
      0  : display <= 7'b1111110;
      1  : display <= 7'b0110000;
      2  : display <= 7'b1101101;
      3  : display <= 7'b1111001;
      4  : display <= 7'b0110011;
      5  : display <= 7'b1011011;
      6  : display <= 7'b1011111;
      7  : display <= 7'b1110000;
      8  : display <= 7'b1111111;
      9  : display <= 7'b1111011;
      10 : display <= 7'b1110111;
      11 : display <= 7'b0011111;
      12 : display <= 7'b1001110;
      13 : display <= 7'b0111101;
      14 : display <= 7'b1001111;
      15 : display <= 7'b1000111;
    endcase
  end
endmodule